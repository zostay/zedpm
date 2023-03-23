package master

import (
	"context"
	"sort"

	"github.com/zostay/zedpm/format"
	"github.com/zostay/zedpm/pkg/storage"
	"github.com/zostay/zedpm/plugin"
)

// Verifies that master.Task implements plugin.Task.
var _ plugin.Task = &Task{}

// taskInfo tracks the plugin name to associate with the execution of a
// particular task.
type taskInfo struct {
	pluginName string
	iface      plugin.Interface
	task       plugin.Task
}

// newTaskInfo constructs a taskInfo object.
func newTaskInfo(
	pluginName string,
	pluginIface plugin.Interface,
	pluginTask plugin.Task,
) *taskInfo {
	return &taskInfo{pluginName, pluginIface, pluginTask}
}

// Task implements plugin.Task by running the operations associated with a set
// of plugins whenever those operations are executed on this task object.
type Task struct {
	taskName string       // name of the task being executed
	ti       *Interface   // link back to the parent interface
	ts       plugin.Tasks // cached list of plugins to work with when executing operations
	taskInfo []*taskInfo  // the plugins associated with each task
}

// newTask constructs a Task.
func newTask(taskName string, ti *Interface, taskInfo []*taskInfo) *Task {
	return &Task{
		taskName: taskName,
		ti:       ti,
		taskInfo: taskInfo,
	}
}

// tasks converts the taskInfo list into a plugin.Tasks, caches the value, and
// returns it. Subsequent calls will use the cached value.
func (t *Task) tasks() plugin.Tasks {
	if t.ts != nil {
		return t.ts
	}

	t.ts = make(plugin.Tasks, len(t.taskInfo))
	for i := range t.taskInfo {
		t.ts[i] = t.taskInfo[i].task
	}
	return t.ts
}

// Setup executes the Setup operation on all associated plugins concurrently.
func (t *Task) Setup(ctx context.Context) error {
	return t.executeTaskOperation(ctx, t.taskName, executeSetup)
}

// Check executes the Check operation on all associated plugins concurrently.
func (t *Task) Check(ctx context.Context) error {
	return t.executeTaskOperation(ctx, t.taskName, executeCheck)
}

// Begin collects all the prioritized operations for the Begin stage of all
// associated plugins and returns a set of master.Operation objects that can
// execute them.
func (t *Task) Begin(ctx context.Context) (plugin.Operations, error) {
	return t.prepareOperations(ctx, plugin.Task.Begin)
}

// Run collects all the prioritized operations for the Run stage of all
// associated plugins and returns a set of master.Operation objects that can
// execute them.
func (t *Task) Run(ctx context.Context) (plugin.Operations, error) {
	return t.prepareOperations(ctx, plugin.Task.Run)
}

// End collects all the prioritized operations for the End stage of all
// associated plugins and returns a set of master.Operation objects that can
// execute them.
func (t *Task) End(ctx context.Context) (plugin.Operations, error) {
	return t.prepareOperations(ctx, plugin.Task.End)
}

// Finish executes the Finish operation on all associated plugins concurrently.
func (t *Task) Finish(ctx context.Context) error {
	return t.executeTaskOperation(ctx, t.taskName, executeFinish)
}

// Teardown executes the Teardown operation on all associated plugins
// concurrently.
func (t *Task) Teardown(ctx context.Context) error {
	return t.executeTaskOperation(ctx, t.taskName, executeTeardown)
}

type taskOperationFunc func(plugin.Task, context.Context) error

// executeBasicStage is a helper for setting up the non-prioritized operation
// executors.
func executeBasicStage(
	stage string,
	opFunc taskOperationFunc,
) func(context.Context, plugin.Task) error {
	return func(ctx context.Context, t plugin.Task) error {
		err := opFunc(t, ctx)
		if err != nil {
			return format.WrapErr(err, "failed %s stage", stage)
		}
		return nil
	}
}

var (
	executeSetup    = executeBasicStage("setup", plugin.Task.Setup)
	executeCheck    = executeBasicStage("check", plugin.Task.Check)
	executeFinish   = executeBasicStage("finish", plugin.Task.Finish)
	executeTeardown = executeBasicStage("teardown", plugin.Task.Teardown)
)

// executeTaskOperation concurrently executes the operation against each plugin
// and returns an error if any occur during execution.
func (t *Task) executeTaskOperation(
	ctx context.Context,
	taskName string,
	op func(context.Context, plugin.Task) error,
) error {
	opfs := make([]plugin.OperationFunc, 0, len(t.taskInfo))
	for i := range t.taskInfo {
		info := t.taskInfo[i]
		opfs = append(opfs, func(ctx context.Context) error {
			ctx, pctx, err := t.ti.ctxFor(ctx, taskName, info.pluginName)
			if err != nil {
				return format.WrapErr(err, "unable to setup plugin context", err)
			}

			err = op(ctx, info.task)
			if err != nil {
				return err
			}

			pctx.UpdateStorage(pctx.StorageChanges())

			return nil
		})
	}

	return executeOperationFuncs(ctx, opfs)
}

// executeOperationFuncs is a generic operation executor function that
// concurrently executes operation functions concurrently, whether being simple
// operations on a task (i.e., Setup, Check, Finish, and Teardown) or the
// prioritized operations (i.e., Begin, Run, and End).
func executeOperationFuncs(
	ctx context.Context,
	opfs []plugin.OperationFunc,
) error {
	return RunTasksAndAccumulateErrors[int, plugin.OperationFunc](
		ctx,
		NewSliceIterator[plugin.OperationFunc](opfs),
		func(ctx context.Context, _ int, opf plugin.OperationFunc) error {
			return opf.Call(ctx)
		})
}

// evaluateOperation performs the preparatory steps for running a prioritized
// operation (i.e., Begin, Run, and End). It collects the information required
// to prepare to execute these operations.
func (t *Task) evaluateOperations(
	ctx context.Context,
	op func(plugin.Task, context.Context) (plugin.Operations, error),
) ([]*operationInfo, error) {
	opInfo := make([]*operationInfo, 0, len(t.taskInfo))
	accChanges := storage.New()
	for _, tInfo := range t.taskInfo {
		ctx, pctx, err := t.ti.ctxFor(ctx, t.taskName, tInfo.pluginName)
		if err != nil {
			return nil, format.WrapErr(err, "unable to setup plugin context", err)
		}

		theseOps, err := op(tInfo.task, ctx)
		if err != nil {
			return nil, err
		}

		accChanges.UpdateStrings(pctx.StorageChanges())

		for _, thisOp := range theseOps {
			info := newOperationInfo(tInfo.pluginName, thisOp)
			opInfo = append(opInfo, info)
		}
	}

	t.ti.properties.UpdateStrings(accChanges.AllSettingsStrings())

	sort.Slice(opInfo, operationInfoLess(opInfo))

	return opInfo, nil
}

// gatherOperations takes a list of operations previously evaluated, sanitizes
// the order value provided, sorts the operations based upon the order, and then
// returns a set of operations that can be executed to perform those operations.
// All plugin operations are grouped by order, so this will always return at
// most 101 operation objects.
func (t *Task) gatherOperations(
	opInfo []*operationInfo,
) plugin.Operations {
	ophs := make(plugin.Operations, 0, len(opInfo))
	var lastOrder plugin.Ordering = -1
	var curOp *OperationHandler
	for _, info := range opInfo {
		order := info.op.Order
		if order < 0 {
			order = 0
		} else if order >= 100 {
			order = 99
		}

		if order > lastOrder {
			curOp = newOperationHandler(
				t.taskName,
				t.ti,
				make([]*operationInfo, 0, len(opInfo)),
			)
			ophs = append(ophs, plugin.Operation{
				Order:  order,
				Action: curOp,
			})
		}

		curOp.opInfo = append(curOp.opInfo, info)

		lastOrder = order
	}

	return ophs
}

// prepareOperations evaluates a prioritized operation (i.e., Begin, Run, and
// End) and gathers these operations into groups ready for execution.
func (t *Task) prepareOperations(
	ctx context.Context,
	op func(plugin.Task, context.Context) (plugin.Operations, error),
) (plugin.Operations, error) {
	opInfo, err := t.evaluateOperations(ctx, op)
	if err != nil {
		return nil, err
	}

	return t.gatherOperations(opInfo), nil
}
