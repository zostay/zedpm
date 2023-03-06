package client

import (
	"context"
	"sort"

	"google.golang.org/grpc"

	"github.com/zostay/zedpm/plugin"
	"github.com/zostay/zedpm/plugin/api"
	"github.com/zostay/zedpm/plugin/translate"
)

// Task implements the plugin.Task interface, which makes calls to the gRPC
// service defined by api.TaskExecutionClient.
type Task struct {
	client api.TaskExecutionClient
	ref    *api.Task_Ref
}

// Setup is a no-op since the api.TaskExecutionClient does not provide this
// functionality.
func (t *Task) Setup(
	_ context.Context,
) error {
	return nil
}

// operation is an internal function for execution a task operation.
func (t *Task) operation(
	ctx context.Context,
	op func(context.Context, *api.Task_Operation_Request, ...grpc.CallOption) (*api.Task_Operation_Response, error),
) error {
	res, err := op(ctx, &api.Task_Operation_Request{
		Task:    t.ref,
		Storage: translate.KVToStringMapString(plugin.KV(ctx)),
	})

	if err != nil {
		return err
	}

	plugin.ApplyChanges(ctx, res.GetStorageUpdate())

	return nil
}

// Check calls the ExecuteCheck gRPC service method.
func (t *Task) Check(
	ctx context.Context,
) error {
	return t.operation(ctx, t.client.ExecuteCheck)
}

// operations is an internal helper that calls a Prepare* function and then
// constructs zero or more Operations to implement the paired Execute* gRPC
// service method.
func (t *Task) operations(
	ctx context.Context,
	prepare func(context.Context, *api.Task_Ref, ...grpc.CallOption) (*api.Task_SubStage_Response, error),
	execute func(context.Context, *api.Task_SubStage_Request, ...grpc.CallOption) (*api.Task_Operation_Response, error),
) (plugin.Operations, error) {
	res, err := prepare(ctx, t.ref)
	if err != nil {
		return nil, err
	}

	orders := res.ProvidedOrders
	sort.Slice(orders, func(i, j int) bool { return orders[i] < orders[j] })

	ops := make(plugin.Operations, len(orders))
	for i, order := range orders {
		ops[i] = plugin.Operation{
			Order: plugin.Ordering(order),
			Action: &Operation{
				parent: t,
				call:   execute,
				order:  order,
			},
		}
	}

	return ops, nil
}

// Begin calls the gRPC PrepareBegin method and returns operations that, when
// called, will call the gRPC ExecuteBegin service method.
func (t *Task) Begin(ctx context.Context) (plugin.Operations, error) {
	return t.operations(ctx, t.client.PrepareBegin, t.client.ExecuteBegin)
}

// Run calls the gRPC PrepareRun method and returns operations that, when
// called, will call the gRPC ExecuteRun service method.
func (t *Task) Run(ctx context.Context) (plugin.Operations, error) {
	return t.operations(ctx, t.client.PrepareRun, t.client.ExecuteRun)
}

// End calls the gRPC PrepareEnd method and returns operations that, when
// called, will call the gRPC ExecuteEnd service method.
func (t *Task) End(ctx context.Context) (plugin.Operations, error) {
	return t.operations(ctx, t.client.PrepareEnd, t.client.ExecuteEnd)
}

// Finish calls the ExecuteFinish gRCP service method.
func (t *Task) Finish(
	ctx context.Context,
) error {
	return t.operation(ctx, t.client.ExecuteFinish)
}

// Teardown is a no-op as the gRPC client does not implement this function.
func (t *Task) Teardown(
	_ context.Context,
) error {
	return nil
}
