package master

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-hclog"

	"github.com/zostay/zedpm/config"
	"github.com/zostay/zedpm/format"
	"github.com/zostay/zedpm/pkg/storage"
	"github.com/zostay/zedpm/plugin"
)

// Prove that master.Interface implements plugin.Interface.
var _ plugin.Interface = &Interface{}

// Interface is a plugin.Interface implementation that aids in the execution of
// a set of other plugin.Interface implementations. This is combined with the
// use of InterfaceExecutor to provide a full set of tools for concurrently
// executing a goal or task.
type Interface struct {
	logger     hclog.Logger                // the logger the master interface uses
	cfg        *config.Config              // the configuration to use during execution
	is         map[string]plugin.Interface // the plugins to execute
	targetName string                      // the target to use when choosing configuration
	pctx       *PhaseContext               // the phase context to track state phase-by-phase
}

// NewInterface creates a new Interface object for the given configuration and
// plugins.
func NewInterface(
	logger hclog.Logger,
	cfg *config.Config,
	is map[string]plugin.Interface,
) *Interface {
	return &Interface{logger, cfg, is, "", NewContext(storage.New())}
}

// GetInterface retrieves the plugin.Interface for the named plugin.
func (ti *Interface) GetInterface(name string) plugin.Interface {
	return ti.is[name]
}

// SetTargetName changes the target used to select the configuration used during
// execution.
func (ti *Interface) SetTargetName(name string) {
	ti.targetName = name
}

// Define records a new value to store in the in-memory properties used during
// interface execution.
func (ti *Interface) Define(values map[string]string) {
	ti.pctx.ApplyChanges(values)
}

// Implements calls Implements on all the associated plugins and returns a
// combined list of all the tasks defined by all the plugins. It fails with an
// error if any plugin fails with an error.
func (ti *Interface) Implements(ctx context.Context) ([]plugin.TaskDescription, error) {
	taskDescs := make([]plugin.TaskDescription, 0, 100)
	for pluginName, iface := range ti.is {
		ctx, err := ti.ctxFor(ctx, "", pluginName)
		if err != nil {
			return nil, err
		}

		tds, err := iface.Implements(ctx)
		if err != nil {
			return nil, err
		}

		tdMap := make(map[string]plugin.TaskDescription, len(tds))
		for _, td := range tds {
			if _, alreadySeen := tdMap[td.Name()]; alreadySeen {
				return nil, fmt.Errorf("duplicate task name %q from plugin %q", td.Name(), pluginName)
			}
			tdMap[td.Name()] = td
		}

		taskDescs = append(taskDescs, tds...)
	}

	return taskDescs, nil
}

// implements is used as an internal check prior to executing Prepare on an
// associated plugin to ensure that the plugin actually implements that task
// before we ask it to perform that task. If it does not implement that task, it
// won't be asked to prepare it.
func (ti *Interface) implements(
	ctx context.Context,
	iface plugin.Interface,
	taskName string,
) (bool, error) {
	// TODO Would it be better to cache information gathered from Implements instead of making extra calls to the plugin?
	taskDescs, err := iface.Implements(ctx)
	if err != nil {
		return false, err
	}

	for _, taskDesc := range taskDescs {
		if taskDesc.Name() == taskName {
			return true, nil
		}
	}
	return false, nil
}

// Goal calls Goal for the given goal name on all associated plugins. If no
// plugin provides a plugin.GoalDescription for this goal, then
// plugin.ErrUnsupportedGoal is returned. Otherwise, the first GoalDescription
// received is returned. If multiple plugins describe a goal with the same name,
// the behavior is non-deterministic.
func (ti *Interface) Goal(
	ctx context.Context,
	name string,
) (plugin.GoalDescription, error) {
	results, err := RunTasksAndAccumulate[string, plugin.Interface, plugin.GoalDescription](
		ctx,
		NewMapIterator[string, plugin.Interface](ti.is),
		func(ctx context.Context, _ string, p plugin.Interface) (plugin.GoalDescription, error) {
			goalDesc, err := p.Goal(ctx, name)
			if errors.Is(err, plugin.ErrUnsupportedGoal) {
				return nil, nil
			} else if err != nil {
				return nil, err
			}
			return goalDesc, nil
		},
	)

	if err != nil {
		return nil, err
	}

	for _, result := range results {
		if result != nil {
			return result, nil
		}
	}

	return nil, plugin.ErrUnsupportedGoal
}

// Prepare calls the Prepare method on all plugins which implements the named
// task. This returns a pointer to a master.Task which is able to execute the
// task for all these plugins. If no plugin implements the named task, then this
// method fails with plugin.ErrUnsupportedTask instead.
func (ti *Interface) Prepare(
	ctx context.Context,
	taskName string,
) (plugin.Task, error) {
	results, err := RunTasksAndAccumulate[string, plugin.Interface, *taskInfo](
		ctx,
		NewMapIterator[string, plugin.Interface](ti.is),
		func(ctx context.Context, pluginName string, iface plugin.Interface) (*taskInfo, error) {
			ctx, err := ti.ctxFor(ctx, taskName, pluginName)
			if err != nil {
				return nil, format.WrapErr(err, "unable to setup plugin context")
			}

			mayPrepare, err := ti.implements(ctx, iface, taskName)
			if err != nil {
				return nil, format.WrapErr(err, "plugin %q failed implements check for task %q", pluginName, taskName)
			}

			if mayPrepare {
				t, err := iface.Prepare(ctx, taskName)
				if err != nil {
					if t != nil {
						return newTaskInfo(pluginName, iface, t), err
					}
					return nil, format.WrapErr(err, "plugin %q failed to run task %q", pluginName, taskName)
				}
				return newTaskInfo(pluginName, iface, t), nil
			}

			return nil, nil
		},
	)

	filteredResults := make([]*taskInfo, 0, len(results))
	for _, result := range results {
		if result == nil {
			continue
		}
		filteredResults = append(filteredResults, result)
	}
	results = filteredResults

	if len(results) > 0 {
		return newTask(taskName, ti, results), err
	}

	if err != nil {
		return nil, err
	}

	return nil, plugin.ErrUnsupportedTask
}

// Cancel performs cancellation for task in progress. It works to immediately
// terminate and close out any resources held by all associated plugins.
func (ti *Interface) Cancel(
	ctx context.Context,
	pluginTask plugin.Task,
) error {
	task := pluginTask.(*Task)
	taskName := finalTaskName(ctx)
	return RunTasksAndAccumulateErrors[int, *taskInfo](
		ctx,
		NewSliceIterator[*taskInfo](task.taskInfo),
		func(ctx context.Context, _ int, p *taskInfo) error {
			ctx, err := ti.ctxFor(ctx, taskName, p.pluginName)
			if err != nil {
				return format.WrapErr(err, "unable to setup plugin context during cancel")
			}

			return p.iface.Cancel(ctx, p.task)
		})
}

// Complete performs completion for the task in progress. It frees up resources
// held by the master.Interface as well as telling each plugin to free up any
// resources associated with task execution on their end for all plugins
// associated with this task.
func (ti *Interface) Complete(
	ctx context.Context,
	pluginTask plugin.Task,
) error {
	task := pluginTask.(*Task)
	taskName := finalTaskName(ctx)
	return RunTasksAndAccumulateErrors[int, *taskInfo](
		ctx,
		NewSliceIterator[*taskInfo](task.taskInfo),
		func(ctx context.Context, _ int, p *taskInfo) error {
			ctx, err := ti.ctxFor(ctx, taskName, p.pluginName)
			if err != nil {
				return format.WrapErr(err, "unable to setup plugin context during complete")
			}

			return p.iface.Complete(ctx, p.task)
		})
}

// ctxFor builds a plugin.Context for the current configuration and target and
// the named task and plugin and associates it with the given context.Context.
func (ti *Interface) ctxFor(
	ctx context.Context,
	taskName string,
	pluginName string,
) (context.Context, error) {
	configProps, err := ti.cfg.ToKV(
		storage.New(),
		taskName,
		ti.targetName,
		pluginName,
	)
	if err != nil {
		return nil, err
	}

	ctx = hclog.WithContext(ctx, ti.logger.With("task", taskName))
	return ti.pctx.withPluginTask(ctx, configProps), nil
}
