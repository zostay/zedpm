package master

import (
	"context"

	"github.com/zostay/zedpm/format"
	"github.com/zostay/zedpm/plugin"
)

// operationInfo tracks the plugin name to which an operation belongs.
type operationInfo struct {
	pluginName string
	op         plugin.Operation
}

// newOperationInfo constructs an operationInfo object.
func newOperationInfo(pluginName string, op plugin.Operation) *operationInfo {
	return &operationInfo{pluginName, op}
}

// operationInfoLess provides a tool for sorting slices of operationInfo
// pointers.
func operationInfoLess(opInfo []*operationInfo) func(i, j int) bool {
	return func(i, j int) bool {
		return opInfo[i].op.Order < opInfo[j].op.Order
	}
}

// OperationHandler implements plugin.OperationHandler and is able to execute
// all the operations for all the plugins associated with executing a particular
// task, operation, stage, and priority order.
type OperationHandler struct {
	taskName string
	ti       *Interface
	opInfo   []*operationInfo
}

// newOperationHandler constructs an OperationHandler.
func newOperationHandler(
	taskName string,
	ti *Interface,
	opInfo []*operationInfo,
) *OperationHandler {
	return &OperationHandler{
		taskName: taskName,
		ti:       ti,
		opInfo:   opInfo,
	}
}

// Call concurrently executes this associated operation and order in all plugins
// that can perform it. It initializes a plugin.Context for each and passes the
// associated configuration through to the plugin. Then, it updates the
// temporary properties for the task using the settings set by the plugin.
func (h *OperationHandler) Call(ctx context.Context) error {
	return RunTasksAndAccumulateErrors[int, *operationInfo](
		ctx,
		NewSliceIterator[*operationInfo](h.opInfo),
		func(ctx context.Context, _ int, info *operationInfo) error {
			ctx, err := h.ti.ctxFor(ctx, h.taskName, info.pluginName)
			if err != nil {
				return format.WrapErr(err, "unable to setup plugin context")
			}

			err = info.op.Action.Call(ctx)
			if err != nil {
				return err
			}

			return nil
		})
}
