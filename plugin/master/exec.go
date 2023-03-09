package master

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/go-hclog"

	"github.com/zostay/zedpm/format"
	"github.com/zostay/zedpm/plugin"
)

// TODO This file (and the project, in general) seems to have some confusion of
// langauge that needs to be straigtened out. I appear to use "stage" to refer to both
// an operational stage within a task and to refer to a groupp of tasks run
// together. Which is it? I think maybe I should use "operation" and "group" as
// the terms and skip "stage" altogether, letting stages and phases be synonyms
// to how we move through the process of task execution generally.

// TODO Take advantage of Golang 1.20's Unwrap() functionality in Error.

// Error is returned by many of the processes here. It represents a list of
// errors. Since concurrency is involved with running multiple tasks at once, it
// is quite possible that multiple failures may occur simultaneously. This error
// implementation collects these errors into a super-error.
type Error []error

// Error returns all the errors inside it as a string.
func (e Error) Error() string {
	msgs := make([]string, len(e))
	for i, err := range e {
		msgs[i] = err.Error()
	}
	return strings.Join(msgs, "; ")
}

// InterfaceExecutor is a tool for executing plugin.Interface objects. It must
// be paired with the master.Interface to help perform this task.
//
// These exist as separate objects because of the separation of concerns between
// these two objects. This object is focused on executing all the operations of
// a task in the correct order and then resolve any errors that occur correctly.
type InterfaceExecutor struct {
	m      *Interface
	logger hclog.Logger
}

// NewExecutor creates a new InterfaceExecutor paired with the given Interface.
func NewExecutor(logger hclog.Logger, m *Interface) *InterfaceExecutor {
	return &InterfaceExecutor{m, logger}
}

// SetTargetName is used to update the target name to use when configuring the
// plugin.Context used to execute plugin.Interface.
func (e *InterfaceExecutor) SetTargetName(name string) {
	e.m.SetTargetName(name)
}

// Define is used to set properties from the command-line or other locations to
// be used when running the plugin.Interface.
func (e *InterfaceExecutor) Define(values map[string]string) {
	e.m.Define(values)
}

// tryCancel executes plugin.Interface.Cancel on the object and internally
// handles the situation where the cancel itself also has an error.
func (e *InterfaceExecutor) tryCancel(
	ctx context.Context,
	taskName string,
	task plugin.Task,
	stage string,
) {
	cancelErr := e.m.Cancel(withFinalTaskName(ctx, taskName), task)
	if cancelErr != nil {
		e.logger.Error("failed while canceling task due to error",
			"stage", stage,
			"task", taskName,
			"error", cancelErr)
	}
}

// logFail logs the information related to a task execution failure.
func (e *InterfaceExecutor) logFail(
	ctx context.Context,
	taskName string,
	stage string,
	err error,
) {
	e.logger.Error("task failed", "stage", stage, "task", taskName, "error", err)
}

// prepare is used to run plugin.Interface.Prepare and handle errors as
// appropriate.
func (e *InterfaceExecutor) prepare(
	ctx context.Context,
	taskName string,
) (plugin.Task, error) {
	task, err := e.m.Prepare(ctx, taskName)
	if err != nil {
		if task != nil {
			e.tryCancel(ctx, taskName, task, "Prepare")
		}
		e.logFail(ctx, taskName, "Prepare", err)
		return nil, err
	}
	return task, nil
}

// taskOperation executes one of the operation-style task stages (i.e., Setup,
// Check, Finish, Teardown).
func (e *InterfaceExecutor) taskOperation(
	ctx context.Context,
	taskName string,
	stage string,
	task plugin.Task,
	op func(context.Context) error,
) error {
	err := op(ctx)
	if err != nil {
		e.tryCancel(ctx, taskName, task, stage)
		e.logFail(ctx, taskName, stage, err)
		return err
	}
	return nil
}

// setup executes the setup stage of the plugin.Task.
func (e *InterfaceExecutor) setup(
	ctx context.Context,
	taskName string,
	task plugin.Task,
) error {
	return e.taskOperation(ctx, taskName, "Setup", task, task.Setup)
}

// check executes the check stage of the plugin.Task.
func (e *InterfaceExecutor) check(
	ctx context.Context,
	taskName string,
	task plugin.Task,
) error {
	return e.taskOperation(ctx, taskName, "Check", task, task.Check)
}

// taskPriorityOperation prepares to run the set of plugin.Operations function
// returned by a prioritized stage (i.e., Begin, Run, and End). And then it runs
// the operations returned by that plugin.Task stage method.
func (e *InterfaceExecutor) taskPriorityOperation(
	ctx context.Context,
	taskName string,
	stage string,
	task plugin.Task,
	prepare func(context.Context) (plugin.Operations, error),
) error {
	ops, err := prepare(ctx)
	if err != nil {
		e.tryCancel(ctx, taskName, task, stage)
		e.logFail(ctx, taskName, stage, err)
		return err
	}

	sort.Slice(ops, plugin.OperationLess(ops))
	for _, op := range ops {
		err := op.Action.Call(ctx)
		if err != nil {
			priStage := fmt.Sprintf("%s:%02d", stage, op.Order)
			e.tryCancel(ctx, taskName, task, priStage)
			e.logFail(ctx, taskName, priStage, err)
			return err
		}
	}

	return nil
}

// begin executes the operations for the Begin phase in priority order.
func (e *InterfaceExecutor) begin(
	ctx context.Context,
	taskName string,
	task plugin.Task,
) error {
	return e.taskPriorityOperation(ctx, taskName, "Begin", task, task.Begin)
}

// run executes the operations for the Run phase in priority order.
func (e *InterfaceExecutor) run(
	ctx context.Context,
	taskName string,
	task plugin.Task,
) error {
	return e.taskPriorityOperation(ctx, taskName, "Run", task, task.Run)
}

// end executes the operations in the End phase in priority order.
func (e *InterfaceExecutor) end(
	ctx context.Context,
	taskName string,
	task plugin.Task,
) error {
	return e.taskPriorityOperation(ctx, taskName, "End", task, task.End)
}

// finish executes the Finish stage of plugin.Task.
func (e *InterfaceExecutor) finish(
	ctx context.Context,
	taskName string,
	task plugin.Task,
) error {
	return e.taskOperation(ctx, taskName, "Finish", task, task.Finish)
}

// teardown executes the Teardown stage of plugin.Teardown.
func (e *InterfaceExecutor) teardown(
	ctx context.Context,
	taskName string,
	task plugin.Task,
) error {
	return e.taskOperation(ctx, taskName, "Teardown", task, task.Teardown)
}

// finalTaskNameKey is the key used with withFinalTaskName and finalTaskName.
type finalTaskNameKey struct{}

// withFinalTaskName inserts the take name into the context for later retrieval
// becase we can't pass it directly through the regular interface even though we
// need it for task execution to work properly.
func withFinalTaskName(ctx context.Context, taskName string) context.Context {
	return context.WithValue(ctx, finalTaskNameKey{}, taskName)
}

// finalTaskName returns the previous stored task name.
func finalTaskName(ctx context.Context) string {
	return ctx.Value(finalTaskNameKey{}).(string)
}

// complete executes the plugin.Interface.Complete method.
func (e *InterfaceExecutor) complete(
	ctx context.Context,
	taskName string,
	task plugin.Task,
) error {
	err := e.m.Complete(withFinalTaskName(ctx, taskName), task)
	if err != nil {
		e.logger.Error("failed while completing task due to error",
			"stage", "Complete",
			"task", taskName,
			"error", err)
	}
	return err
}

// ExecuteAllStages sorts all the stages for execution of a task into groups.
// These groups may represent implementations of multiple tasks to achieve a
// goal or a sub-task of a goal. These are each executed concurrently in order
// from most required to least required to complete the goal or subtask in its
// entirety.
//
// Lots of error checking is performed. If anything goes wrong, this operation
// will fail with an error.
func (e *InterfaceExecutor) ExecuteAllStages(
	ctx context.Context,
	group *TaskGroup,
) error {
	stages, err := group.ExecutionGroups()
	if err != nil {
		return fmt.Errorf("failed to break down goal %q into stages: %v", group.Goal.Name(), err)
	}

	for _, stage := range stages {
		err := e.ExecuteStage(ctx, stage)
		if err != nil {
			stageNames := make([]string, len(stage))
			for i, task := range stage {
				stageNames[i] = task.Name()
			}
			return fmt.Errorf("failed to execute stage (%s): %v", format.And(stageNames...), err)
		}
	}

	return nil
}

// ExecuteStage will execute a set of tasks concurrently. Any errors that occur
// executing this stage will be returned as an Error.
func (e *InterfaceExecutor) ExecuteStage(
	ctx context.Context,
	stage []plugin.TaskDescription,
) error {
	return RunTasksAndAccumulateErrors[int, plugin.TaskDescription](ctx,
		NewSliceIterator[plugin.TaskDescription](stage),
		func(ctx context.Context, _ int, task plugin.TaskDescription) error {
			return e.Execute(ctx, task.Name())
		},
	)
}

// Execute will execute a single task and return an error if execution fails.
func (e *InterfaceExecutor) Execute(
	ctx context.Context,
	taskName string,
) error {
	task, err := e.prepare(ctx, taskName)
	if err != nil {
		return fmt.Errorf("failed to prepare task %q: %v", taskName, err)
	}

	stdOps := []struct {
		name     string
		function func(context.Context, string, plugin.Task) error
	}{
		{"setup", e.setup},
		{"check", e.check},
		{"begin", e.begin},
		{"run", e.run},
		{"end", e.end},
		{"finish", e.finish},
		{"teardown", e.teardown},
		{"complete", e.complete},
	}

	for _, stdOp := range stdOps {
		err = stdOp.function(ctx, taskName, task)
		if err != nil {
			return fmt.Errorf("failed to execute operation %s: %v", stdOp.name, err)
		}
	}

	return nil
}

// TaskGroups builds an returns a slice of TaskGroup objects that will be
// executed as part of this InterfaceExecutor.
func (e *InterfaceExecutor) TaskGroups(
	ctx context.Context,
) ([]*TaskGroup, error) {
	tasks, err := e.m.Implements(ctx)
	if err != nil {
		return nil, err
	}

	taskGroupMap := make(map[string]*TaskGroup, len(tasks))

	for _, task := range tasks {
		goalName, err := plugin.GoalName(task)
		if err != nil {
			return nil, err
		}

		if group, groupExists := taskGroupMap[goalName]; groupExists {
			group.Tasks = append(group.Tasks, task)
			continue
		}

		goal, err := e.m.Goal(ctx, goalName)
		if err != nil {
			return nil, err
		}

		taskGroupMap[goalName] = &TaskGroup{
			Tree:  "/" + goal.Name(),
			Goal:  goal,
			Tasks: []plugin.TaskDescription{task},
		}
	}

	out := make([]*TaskGroup, 0, len(tasks))
	for _, group := range taskGroupMap {
		out = append(out, group)
	}

	for _, group := range out {
		sort.Slice(group.Tasks, func(i, j int) bool {
			return group.Tasks[i].Name() < group.Tasks[j].Name()
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Goal.Name() < out[j].Goal.Name()
	})

	return out, nil
}
