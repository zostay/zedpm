package client

import (
	"context"
	"sort"

	"google.golang.org/grpc"

	"github.com/zostay/zedpm/plugin"
	"github.com/zostay/zedpm/plugin/api"
	"github.com/zostay/zedpm/plugin/translate"
)

type Task struct {
	client api.TaskExecutionClient
	ref    *api.Task_Ref
}

func (t *Task) Setup(
	_ context.Context,
) error {
	return nil
}

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

func (t *Task) Check(
	ctx context.Context,
) error {
	return t.operation(ctx, t.client.ExecuteCheck)
}

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

func (t *Task) Begin(ctx context.Context) (plugin.Operations, error) {
	return t.operations(ctx, t.client.PrepareBegin, t.client.ExecuteBegin)
}

func (t *Task) Run(ctx context.Context) (plugin.Operations, error) {
	return t.operations(ctx, t.client.PrepareRun, t.client.ExecuteRun)
}

func (t *Task) End(ctx context.Context) (plugin.Operations, error) {
	return t.operations(ctx, t.client.PrepareEnd, t.client.ExecuteEnd)
}

func (t *Task) Finish(
	ctx context.Context,
) error {
	return t.operation(ctx, t.client.ExecuteFinish)
}

func (t *Task) Teardown(
	_ context.Context,
) error {
	return nil
}
