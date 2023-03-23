package master

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/zostay/zedpm/config"
	"github.com/zostay/zedpm/format"
	"github.com/zostay/zedpm/pkg/group"
	"github.com/zostay/zedpm/plugin"
)

// TODO This file (and the project, in general) seems to have some confusion of
// langauge that needs to be straigtened out. I appear to use "stage" to refer to both
// an operational stage within a task and to refer to a groupp of tasks run
// together. Which is it? I think maybe I should use "operation" and "group" as
// the terms and skip "stage" altogether, letting stages and phases be synonyms
// to how we move through the process of task execution generally.

// TODO Take advantage of Golang 1.20's Unwrap() functionality in Error.

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
		errMsg := format.Err(cancelErr)
		e.logger.Error("failed while canceling task due to error",
			"stage", stage,
			"task", taskName,
			"error", errMsg)
	}
}

// logFail logs the information related to a task execution failure.
func (e *InterfaceExecutor) logFail(
	ctx context.Context,
	taskName string,
	stage string,
	err error,
) {
	e.logger.Error("task failed",
		"stage", stage,
		"task", taskName,
		"error", format.Err(err))
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

// ExecuteGoal executes all the tasks in a goal. Tasks are grouped into phases.
// Each phase is run one at a time in order. These may be executed concurrently.
// The tasks within each phase are run simultaneously and interleaved (with
// individual operations sometimes running concurrently).
func (e *InterfaceExecutor) ExecuteGoal(
	ctx context.Context,
	goal *group.Goal,
) error {
	phases := goal.ExecutionPhases()

	for _, phase := range phases {
		err := e.ExecutePhase(ctx, phase)
		return err
	}

	return nil
}

// ExecutePhase executes all the tasks in a phase. Tasks in a phase are executed
// simultaneously with operations interleaved and run concurrently according to
// operation order.
//
// First, every task is concurrently setup, then checked. Then the operation of
// the begin phase are run in priority order, with any operations having the
// same priority being run concurrently. Then the operations of the run phase
// run, again in priority order. And then the operations of the end phase are
// run in priority order. Finally, the tasks are finished and torn down
// concurrently.
func (e *InterfaceExecutor) ExecutePhase(
	ctx context.Context,
	phase *group.Phase,
) error {
	ops := []OperationExecutor{
		&SimpleExecutor{"setup", e, phase.Tasks(), plugin.Task.Setup},
		&SimpleExecutor{"check", e, phase.Tasks(), plugin.Task.Check},

		&StagedExecutor{"begin", e, phase.Tasks(), plugin.Task.Begin},
		&StagedExecutor{"run", e, phase.Tasks(), plugin.Task.Run},
		&StagedExecutor{"end", e, phase.Tasks(), plugin.Task.End},

		&SimpleExecutor{"finish", e, phase.Tasks(), plugin.Task.Finish},
		&SimpleExecutor{"teardown", e, phase.Tasks(), plugin.Task.Teardown},

		&CompletionExecutor{e, phase.Tasks()},
	}

	for _, op := range ops {
		err := op.Execute(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// PotentialGoalsPhasesAndTasks builds an returns a slice of TaskGroup objects that will be
// executed as part of this InterfaceExecutor.
func (e *InterfaceExecutor) PotentialGoalsPhasesAndTasks(
	ctx context.Context,
) ([]*group.Goal, error) {
	tasks, err := e.m.Implements(ctx)
	if err != nil {
		return nil, err
	}

	goalMap := make(map[string]plugin.GoalDescription, len(tasks))

	for _, task := range tasks {
		goalName, _, _, err := config.GoalPhaseAndTaskName(task.Name())
		if err != nil {
			return nil, err
		}

		if goal := goalMap[goalName]; goal != nil {
			continue
		}

		goal, err := e.m.Goal(ctx, goalName)
		if err != nil {
			return nil, err
		}

		goalMap[goalName] = goal
	}

	return group.SetupGroups(tasks, goalMap)
}
