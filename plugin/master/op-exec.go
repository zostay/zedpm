package master

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"

	"github.com/zostay/zedpm/format"
	"github.com/zostay/zedpm/pkg/group"
	"github.com/zostay/zedpm/plugin"
)

// OperationExecutor is the abstraction that abstracts away the differences
// between the different types of operations that are executed in the course of
// executing an operation.
type OperationExecutor interface {
	Execute(context.Context) error
}

// OperationFunc is a more abstract version of plugin.OperationFunc.
type OperationFunc func(plugin.Task, context.Context) error

// PrepareFunc is an abstract version of plugin.Task operation functuons.
type PrepareFunc func(plugin.Task, context.Context) (plugin.Operations, error)

// SimpleExecutor is an executor that can execute one of the simple operations
// on a task. As of this writing, this includes the setup, check, finish, and
// teardown operations.
type SimpleExecutor struct {
	stageName string
	exec      *InterfaceExecutor
	tasks     []plugin.TaskDescription
	run       OperationFunc
}

// Execute runs all the configured operation stage for all tasks concurrently,
// collects the errors, and returns them.
func (s *SimpleExecutor) Execute(
	ctx context.Context,
) error {
	return RunTasksAndAccumulateErrors[int, plugin.TaskDescription](ctx,
		NewSliceIterator[plugin.TaskDescription](s.tasks),
		func(ctx context.Context, _ int, taskDesc plugin.TaskDescription) error {
			taskName := taskDesc.Name()

			logger := hclog.FromContext(ctx)
			logger = logger.With(
				"@operation", s.stageName,
				"@task", taskName,
			)
			ctx = hclog.WithContext(ctx, logger)

			task, err := s.exec.prepare(ctx, taskName)
			if err != nil {
				return format.WrapErr(err, "failed to prepare task %q", taskName)
			}

			err = s.run(task, ctx)
			if err != nil {
				return format.WrapErr(err, "failed to execute operation %s", s.stageName)
			}

			return nil
		},
	)
}

// CompletionExecutor executes the Complete operation that must run when a task
// is complete. This runs Complete for all tasks concurrently, collects the
// errors, and returns them.
type CompletionExecutor struct {
	exec  *InterfaceExecutor
	tasks []plugin.TaskDescription
	phase *group.Phase
}

func (c *CompletionExecutor) Execute(
	ctx context.Context,
) error {
	return RunTasksAndAccumulateErrors[int, plugin.TaskDescription](ctx,
		NewSliceIterator[plugin.TaskDescription](c.tasks),
		func(ctx context.Context, _ int, taskDesc plugin.TaskDescription) error {
			taskName := taskDesc.Name()
			task, err := c.exec.prepare(ctx, taskName)
			if err != nil {
				c.exec.logger.Error("unknown error while completing task",
					"@operation", "Complete",
					"@task", taskName,
					"error", format.Err(err))
				return err
			}

			ctx = withFinalTaskName(ctx, taskName)
			err = c.exec.m.Complete(ctx, task)
			if err != nil {
				c.exec.logger.Error("failed while completing task due to error",
					"@operation", "Complete",
					"@task", taskName,
					"error", format.Err(err))
			}

			return err
		},
	)
}

// StagedExecutor handles the primary staged phases: begin, run, and end. This
// works by pulling all the prioritized operations and then executing each group
// in priority order. Operations with the same priority will be executed
// concurrently. Otherwise, operations are ordered by their plugin.Ordering,
// from lowest numbered ordering to highest. If an error happens at any point,
// the errors are returned and the operation terminates immediately without
// continuing on to any remaining operations.
type StagedExecutor struct {
	stageName string
	exec      *InterfaceExecutor
	tasks     []plugin.TaskDescription
	prepare   PrepareFunc
}

type opInfo struct {
	op       plugin.Operation
	task     plugin.Task
	taskName string
}

// prepareTasks returns a slice of prepared plugin.Task objects in the same
// order as the plugin.TaskDescription objects in the tasks field.
func (s *StagedExecutor) prepareTasks(
	ctx context.Context,
) ([]plugin.Task, error) {
	tasks := make([]plugin.Task, 0, len(s.tasks))
	for _, taskDesc := range s.tasks {
		taskName := taskDesc.Name()

		logger := hclog.FromContext(ctx)
		logger = logger.With(
			"@task", taskName,
			"@operation", s.stageName,
		)
		ctx := hclog.WithContext(ctx, logger)

		task, err := s.exec.prepare(ctx, taskName)
		if err != nil {
			err = fmt.Errorf("failed to prepare task %q: %w", taskName, err)
			s.exec.tryCancel(ctx, taskName, task, s.stageName)
			s.exec.logFail(taskName, s.stageName, err)
			return nil, err
		}

		tasks = append(tasks, task)
	}
	return tasks, nil
}

// collateTasks returns the collected operations gathered together in execution
// order with operations that need to be executed concurrently being grouped
// together.
func (s *StagedExecutor) collateTasks(
	ctx context.Context,
	tasks []plugin.Task,
) (map[plugin.Ordering][]opInfo, error) {
	ops := make(map[plugin.Ordering][]opInfo, 100)
	for i, task := range tasks {
		taskName := s.tasks[i].Name()

		logger := hclog.FromContext(ctx)
		logger = logger.With(
			"@task", taskName,
			"@operation", s.stageName,
		)
		ctx := hclog.WithContext(ctx, logger)

		theseOps, err := s.prepare(task, ctx)
		if err != nil {
			err = fmt.Errorf("failed to prepare task %q: %w", taskName, err)
			s.exec.tryCancel(ctx, taskName, task, s.stageName)
			s.exec.logFail(taskName, s.stageName, err)
			return nil, err
		}

		for _, thisOp := range theseOps {
			if _, listExists := ops[thisOp.Order]; !listExists {
				ops[thisOp.Order] = make([]opInfo, 0, 1)
			}

			ops[thisOp.Order] = append(ops[thisOp.Order], opInfo{
				op:       thisOp,
				task:     task,
				taskName: taskName,
			})
		}
	}

	return ops, nil
}

// executeOperations executes all the operations, each in their groups.
func (s *StagedExecutor) executeOperations(
	ctx context.Context,
	ops map[plugin.Ordering][]opInfo,
) error {
	for i := plugin.Ordering(0); i < 100; i++ {
		if stageOps, opsExist := ops[i]; opsExist {
			err := RunTasksAndAccumulateErrors[int, opInfo](ctx,
				NewSliceIterator[opInfo](stageOps),
				func(ctx context.Context, _ int, opInfo opInfo) error {
					taskName := opInfo.taskName
					priStage := fmt.Sprintf("%s:%v", s.stageName, opInfo.op.Order)

					logger := hclog.FromContext(ctx)
					logger.With(
						"@task", taskName,
						"@operation", priStage,
					)
					ctx = hclog.WithContext(ctx, logger)

					err := opInfo.op.Action.Call(ctx)
					if err != nil {
						err = fmt.Errorf("failed while executing stage %s of task %q: %w", priStage, taskName, err)
						s.exec.tryCancel(ctx, taskName, opInfo.task, priStage)
						s.exec.logFail(taskName, priStage, err)
						return err
					}
					return nil
				},
			)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Execute builds a plan for executing a staged operation and then executes that
// plan. Any errors that occur along the way are returned.
func (s *StagedExecutor) Execute(
	ctx context.Context,
) error {
	tasks, err := s.prepareTasks(ctx)
	if err != nil {
		return err
	}

	ops, err := s.collateTasks(ctx, tasks)
	if err != nil {
		return err
	}

	return s.executeOperations(ctx, ops)
}
