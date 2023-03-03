package client

import (
	"context"

	"google.golang.org/grpc"

	"github.com/zostay/zedpm/plugin"
	"github.com/zostay/zedpm/plugin/api"
	"github.com/zostay/zedpm/plugin/translate"
)

type Operation struct {
	parent *Task
	call   func(context.Context, *api.Task_SubStage_Request, ...grpc.CallOption) (*api.Task_Operation_Response, error)
	order  int32
}

func (o *Operation) Call(ctx context.Context) error {
	res, err := o.call(ctx, &api.Task_SubStage_Request{
		Request: &api.Task_Operation_Request{
			Task:    o.parent.ref,
			Storage: translate.KVToStringMapString(plugin.KV(ctx)),
		},
		SubStage: o.order,
	})

	if err != nil {
		return err
	}

	plugin.ApplyChanges(ctx, res.GetStorageUpdate())

	return nil
}
