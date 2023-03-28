package client

import (
	"context"

	"google.golang.org/grpc"

	"github.com/zostay/zedpm/plugin/api"
	"github.com/zostay/zedpm/plugin/translate"
)

// Operation implements plugin.OperationHandler, which makes calls to the
// api.TaskExecutionClient gRPC client.
type Operation struct {
	parent *Task
	call   func(context.Context, *api.Task_SubStage_Request, ...grpc.CallOption) (*api.Task_Operation_Response, error)
	order  int32
}

// Call executes the operaiton by calling one of the Execute* gRPC service
// methods.
func (o *Operation) Call(ctx context.Context) error {
	res, err := o.call(ctx, &api.Task_SubStage_Request{
		Request: &api.Task_Operation_Request{
			Task:       o.parent.ref,
			Storage:    translate.KVToStringMapString(KV(ctx)),
			AddedFiles: ListAdded(ctx),
		},
		SubStage: o.order,
	})

	if err != nil {
		return err
	}

	ApplyChanges(ctx, res.GetStorageUpdate())
	ToAdd(ctx, res.GetAddedFiles())

	return nil
}
