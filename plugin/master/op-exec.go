package master

import (
	"context"
	"fmt"

	"github.com/zostay/zedpm/format"
	"github.com/zostay/zedpm/plugin"
)

type OperationExecutor interface {
	Execute(context.Context) error
}

type OperationFunc func(plugin.Task, context.Context) error
type PrepareFunc func(plugin.Task, context.Context) (plugin.Operations, error)

type SimpleExecutor struct {
	stageName string
	exec      *InterfaceExecutor
	tasks     []plugin.TaskDescription
	run       OperationFunc
}

func (s *SimpleExecutor) Execute(
	ctx context.Context,
) error {
	return RunTasksAndAccumulateErrors[int, plugin.TaskDescription](ctx,
		NewSliceIterator[plugin.TaskDescription](s.tasks),
		func(ctx context.Context, _ int, taskDesc plugin.TaskDescription) error {
			taskName := taskDesc.Name()
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

type CompletionExecutor struct {
	exec  *InterfaceExecutor
	tasks []plugin.TaskDescription
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
					"stage", "Complete",
					"task", taskName,
					"error", format.Err(err))
				return err
			}

			ctx = withFinalTaskName(ctx, taskName)
			err = c.exec.m.Complete(ctx, task)
			if err != nil {
				c.exec.logger.Error("failed while completing task due to error",
					"stage", "Complete",
					"task", taskName,
					"error", format.Err(err))
			}

			return err
		},
	)
}

type StagedExecutor struct {
	stageName string
	exec      *InterfaceExecutor
	tasks     []plugin.TaskDescription
	prepare   PrepareFunc
}

func (s *StagedExecutor) Execute(
	ctx context.Context,
) error {
	type opInfo struct {
		op       plugin.Operation
		task     plugin.Task
		taskName string
	}

	tasks := make([]plugin.Task, len(s.tasks))
	for _, taskDesc := range s.tasks {
		taskName := taskDesc.Name()
		task, err := s.exec.prepare(ctx, taskName)
		if err != nil {
			err = fmt.Errorf("failed to prepare task %q: %w", taskName, err)
			s.exec.tryCancel(ctx, taskName, task, s.stageName)
			s.exec.logFail(ctx, taskName, s.stageName, err)
			return err
		}

		tasks = append(tasks, task)
	}

	ops := make(map[plugin.Ordering][]opInfo, 100)
	for i, task := range tasks {
		taskName := s.tasks[i].Name()
		theseOps, err := s.prepare(task, ctx)
		if err != nil {
			err = fmt.Errorf("failed to prepare task %q: %w", taskName, err)
			s.exec.tryCancel(ctx, taskName, task, s.stageName)
			s.exec.logFail(ctx, taskName, s.stageName, err)
			return err
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

	for i := plugin.Ordering(0); i < 100; i++ {
		if stageOps, opsExist := ops[i]; opsExist {
			err := RunTasksAndAccumulateErrors[int, opInfo](ctx,
				NewSliceIterator[opInfo](stageOps),
				func(ctx context.Context, _ int, opInfo opInfo) error {
					err := opInfo.op.Action.Call(ctx)
					if err != nil {
						priStage := fmt.Sprintf("%s:%v", s.stageName, opInfo.op.Order)
						taskName := opInfo.taskName
						err = fmt.Errorf("failed while executing tsage %s of task %q: %w", priStage, taskName, err)
						s.exec.tryCancel(ctx, taskName, opInfo.task, priStage)
						s.exec.logFail(ctx, taskName, priStage, err)
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
