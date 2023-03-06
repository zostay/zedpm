package client

import (
	"context"
	"strings"

	"github.com/zostay/zedpm/plugin"
	"github.com/zostay/zedpm/plugin/api"
	"github.com/zostay/zedpm/plugin/translate"
)

// Verifies that Interface is a plugin.Interface.
var _ plugin.Interface = &Interface{}

// Interface implements plugin.Interface to map calls to that interface onto
// calls to the api.TaskExecutionClient.
type Interface struct {
	client api.TaskExecutionClient
}

// NewGRPCTaskInterface returns a new plugin.Interface wrapper around
// api.TaskExecutionClient.
func NewGRPCTaskInterface(client api.TaskExecutionClient) *Interface {
	return &Interface{client}
}

// Implements calls the Implements gRPC service method.
func (c *Interface) Implements(
	ctx context.Context,
) ([]plugin.TaskDescription, error) {
	res, err := c.client.Implements(ctx, &api.Task_Implements_Request{})
	if err != nil {
		return nil, err
	}
	return translate.APITaskDescriptorsToPluginTaskDescriptions(res.GetTasks()), nil
}

// Goal calls the Goal gRPC service method.
func (c *Interface) Goal(
	ctx context.Context,
	goalName string,
) (plugin.GoalDescription, error) {
	res, err := c.client.Goal(ctx, &api.Task_Goal_Request{
		Name: goalName,
	})
	if err != nil {
		if strings.Contains(err.Error(), plugin.ErrUnsupportedGoal.Error()) {
			return nil, plugin.ErrUnsupportedGoal
		}
		return nil, err
	}
	return translate.APIGoalDescriptorToPluginGoalDescription(res.GetDefinition()), nil
}

// Prepare calls the Prepare gRPC service method.
func (c *Interface) Prepare(
	ctx context.Context,
	taskName string,
) (plugin.Task, error) {
	res, err := c.client.Prepare(ctx,
		&api.Task_Prepare_Request{
			Name:         taskName,
			GlobalConfig: translate.KVToAPIConfig(plugin.KV(ctx)),
		},
	)
	if err != nil {
		if strings.Contains(err.Error(), plugin.ErrUnsupportedTask.Error()) {
			return nil, plugin.ErrUnsupportedTask
		}
		return nil, err
	}

	plugin.ApplyChanges(ctx, res.GetStorage())

	return &Task{
		client: c.client,
		ref:    res.GetTask(),
	}, nil
}

// Cancel calls the Cancel gRPC srvice method.
func (c *Interface) Cancel(
	ctx context.Context,
	task plugin.Task,
) error {
	ref := task.(*Task).ref
	_, err := c.client.Cancel(ctx, &api.Task_Cancel_Request{
		Task:    ref,
		Storage: translate.KVToStringMapString(plugin.KV(ctx)),
	})
	return err
}

// Complete calls the Complete gRPC method.
func (c *Interface) Complete(
	ctx context.Context,
	task plugin.Task,
) error {
	ref := task.(*Task).ref
	_, err := c.client.Complete(ctx, &api.Task_Complete_Request{
		Task:    ref,
		Storage: translate.KVToStringMapString(plugin.KV(ctx)),
	})
	return err
}
