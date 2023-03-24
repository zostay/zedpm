package service

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zostay/zedpm/pkg/storage"
	"github.com/zostay/zedpm/plugin"
	"github.com/zostay/zedpm/plugin/api"
	"github.com/zostay/zedpm/plugin/translate"
)

// Verifies that TaskExecution implements api.TaskExecutionServer.
var _ api.TaskExecutionServer = &TaskExecution{}

// TaskState is used to track the task state for a running plugin task.
type TaskState struct {
	// Task is the plugin.Task object that is being executed for this task state.
	Task plugin.Task

	// Context is the plugin.Context used when executing this task.
	Context *plugin.Context
}

// TaskExecution implements api.TaskExecutionServer and maps incoming gRPC
// service calls to the plugin.Interface.
type TaskExecution struct {
	api.UnimplementedTaskExecutionServer

	Impl   plugin.Interface
	logger hclog.Logger
	state  map[string]map[string]*TaskState
}

// NewGRPCTaskExecution returns a new TaskExecution object that will map
// incoming gRPC service calls to a plugin.Interface implementation.
func NewGRPCTaskExecution(
	logger hclog.Logger,
	impl plugin.Interface,
) *TaskExecution {
	taskDescs, err := impl.Implements(context.Background())
	if err != nil {
		return nil
	}

	state := make(map[string]map[string]*TaskState, len(taskDescs))
	for _, taskDesc := range taskDescs {
		state[taskDesc.Name()] = make(map[string]*TaskState, 1)
	}

	return &TaskExecution{
		Impl:   impl,
		logger: logger,
		state:  state,
	}
}

// generateStateId is the internal function used to generate state_id returned
// as part of a task reference.
func generateStateId() string {
	return ulid.Make().String()
}

// Implements maps the gRPC Implements service method to the Implements method
// of plugin.Interface.
func (s *TaskExecution) Implements(
	ctx context.Context,
	_ *api.Task_Implements_Request,
) (*api.Task_Implements_Response, error) {
	taskDescs, err := s.Impl.Implements(ctx)
	if err != nil {
		return nil, err
	}
	return &api.Task_Implements_Response{
		Tasks: translate.PluginTaskDescriptionsToAPITaskDescriptors(taskDescs),
	}, nil
}

// Goal maps the gRPC Goal service method to the Goal method of
// plugin.Interface.
func (s *TaskExecution) Goal(
	ctx context.Context,
	request *api.Task_Goal_Request,
) (*api.Task_Goal_Response, error) {
	pctx := plugin.NewContext(s.logger, storage.New())
	ctx = plugin.InitializeContext(ctx, pctx)

	goalDesc, err := s.Impl.Goal(ctx, request.GetName())
	if err != nil {
		return nil, err
	}

	return &api.Task_Goal_Response{
		Definition: translate.PluginGoalDescriptionToAPIGoalDescriptor(goalDesc),
	}, nil
}

// Prepare maps the gRPC Prepare service method to the Prepare method of
// plugin.Interface.
func (s *TaskExecution) Prepare(
	ctx context.Context,
	request *api.Task_Prepare_Request,
) (*api.Task_Prepare_Response, error) {
	globalConfig := request.GetGlobalConfig()

	kv := translate.APIConfigToKV(globalConfig)
	pctx := plugin.NewContext(s.logger, kv)
	ctx = plugin.InitializeContext(ctx, pctx)

	task, err := s.Impl.Prepare(ctx, request.GetName())
	if err != nil {
		return nil, err
	}

	state := &TaskState{
		Task:    task,
		Context: pctx,
	}

	name := request.GetName()
	id := generateStateId()
	s.state[name][id] = state

	res, err := s.executeStage(ctx, &api.Task_Operation_Request{
		Task: &api.Task_Ref{
			Name:    name,
			StateId: id,
		},
		Storage: map[string]string{},
	}, plugin.Task.Setup)
	if err != nil {
		return nil, err
	}

	return &api.Task_Prepare_Response{
		Task: &api.Task_Ref{
			Name:    request.GetName(),
			StateId: id,
		},
		Storage: res.GetStorageUpdate(),
	}, nil
}

// deref takes an incoming task reference and turns it into a TaskState object.
func (s *TaskExecution) deref(ref *api.Task_Ref) (*TaskState, error) {
	name := ref.GetName()
	id := ref.GetStateId()
	task := s.state[name][id]
	if task == nil {
		return nil, status.Errorf(codes.NotFound, "the task named %q with state ID %q not found", name, id)
	}

	return task, nil
}

// closeTask performs final operations on a task, performs cleanup, and deletes
// the task reference from the internal state cache.
func (s *TaskExecution) closeTask(
	ctx context.Context,
	taskRef *api.Task_Ref,
	storage map[string]string,
	completed bool,
) error {
	_, err := s.deref(taskRef)
	if err != nil {
		return err
	}

	_, err = s.executeStage(ctx, &api.Task_Operation_Request{
		Task:    taskRef,
		Storage: storage,
	}, plugin.Task.Teardown)

	state, derefErr := s.deref(taskRef)
	if derefErr != nil {
		plugin.Logger(ctx).Error("fatal error during cancel: %v", derefErr)
		// TODO This here is a sign that there's a problem with this API layout. This should be fixed.
		panic("fatal error during plugin cancellation")
	}

	ctx = plugin.InitializeContext(ctx, state.Context)

	if err != nil { //nolint:gocritic // prefer this as an if-else
		anotherErr := s.Impl.Cancel(ctx, state.Task)
		plugin.Logger(ctx).Error("error during plugin cancel: %v", anotherErr)
	} else if !completed {
		err = s.Impl.Cancel(ctx, state.Task)
	} else {
		err = s.Impl.Complete(ctx, state.Task)
	}

	delete(s.state[taskRef.GetName()], taskRef.GetStateId())

	return err
}

// Cancel implements the Cancel gRPC service method and performs final cleanup
// necessary for cancelling a task.
func (s *TaskExecution) Cancel(
	ctx context.Context,
	request *api.Task_Cancel_Request,
) (*api.Task_Cancel_Response, error) {
	err := s.closeTask(ctx, request.GetTask(), request.GetStorage(), false)
	if err != nil {
		return nil, err
	}
	return &api.Task_Cancel_Response{}, nil
}

// Complete implements the Complete gRPC service method and performs final
// cleanup necessary for completing a task.
func (s *TaskExecution) Complete(
	ctx context.Context,
	request *api.Task_Complete_Request,
) (*api.Task_Complete_Response, error) {
	err := s.closeTask(ctx, request.GetTask(), request.GetStorage(), true)
	if err != nil {
		return nil, err
	}
	return &api.Task_Complete_Response{}, nil
}
