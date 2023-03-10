// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: task-interface.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TaskExecutionClient is the client API for TaskExecution service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TaskExecutionClient interface {
	// Implements maps onto the plugin.Interface.Implements method.
	Implements(ctx context.Context, in *Task_Implements_Request, opts ...grpc.CallOption) (*Task_Implements_Response, error)
	// Goal maps onto the plugin.Interface.Goal method.
	Goal(ctx context.Context, in *Task_Goal_Request, opts ...grpc.CallOption) (*Task_Goal_Response, error)
	// Prepare maps onto the plugin.Interface.Prepare method.
	Prepare(ctx context.Context, in *Task_Prepare_Request, opts ...grpc.CallOption) (*Task_Prepare_Response, error)
	// Cancel maps onto the plugin.Interface.Cancel method.
	Cancel(ctx context.Context, in *Task_Cancel_Request, opts ...grpc.CallOption) (*Task_Cancel_Response, error)
	// Complete maps onto the plugin.Interface.Complete method.
	Complete(ctx context.Context, in *Task_Complete_Request, opts ...grpc.CallOption) (*Task_Complete_Response, error)
	// ExecuteCheck maps onto the plugin.Task.Check method.
	ExecuteCheck(ctx context.Context, in *Task_Operation_Request, opts ...grpc.CallOption) (*Task_Operation_Response, error)
	// PrepareBegin maps onto the plugin.Task.Begin method.
	PrepareBegin(ctx context.Context, in *Task_Ref, opts ...grpc.CallOption) (*Task_SubStage_Response, error)
	// ExecuteBegin executes operations related to the plugin.Task.Begin method.
	ExecuteBegin(ctx context.Context, in *Task_SubStage_Request, opts ...grpc.CallOption) (*Task_Operation_Response, error)
	// PrepareRun maps onto the plugin.Task.Run method.
	PrepareRun(ctx context.Context, in *Task_Ref, opts ...grpc.CallOption) (*Task_SubStage_Response, error)
	// ExecuteRun executes operations related to the plugin.Task.Run method.
	ExecuteRun(ctx context.Context, in *Task_SubStage_Request, opts ...grpc.CallOption) (*Task_Operation_Response, error)
	// PrepareEnd maps onto the plugin.Task.End method.
	PrepareEnd(ctx context.Context, in *Task_Ref, opts ...grpc.CallOption) (*Task_SubStage_Response, error)
	// ExecuteEnd executes operations related to the plugin.Task.End method.
	ExecuteEnd(ctx context.Context, in *Task_SubStage_Request, opts ...grpc.CallOption) (*Task_Operation_Response, error)
	// ExecuteFinish maps onto the plugin.Task.Finish method.
	ExecuteFinish(ctx context.Context, in *Task_Operation_Request, opts ...grpc.CallOption) (*Task_Operation_Response, error)
}

type taskExecutionClient struct {
	cc grpc.ClientConnInterface
}

func NewTaskExecutionClient(cc grpc.ClientConnInterface) TaskExecutionClient {
	return &taskExecutionClient{cc}
}

func (c *taskExecutionClient) Implements(ctx context.Context, in *Task_Implements_Request, opts ...grpc.CallOption) (*Task_Implements_Response, error) {
	out := new(Task_Implements_Response)
	err := c.cc.Invoke(ctx, "/zedpm.plugin.TaskExecution/Implements", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskExecutionClient) Goal(ctx context.Context, in *Task_Goal_Request, opts ...grpc.CallOption) (*Task_Goal_Response, error) {
	out := new(Task_Goal_Response)
	err := c.cc.Invoke(ctx, "/zedpm.plugin.TaskExecution/Goal", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskExecutionClient) Prepare(ctx context.Context, in *Task_Prepare_Request, opts ...grpc.CallOption) (*Task_Prepare_Response, error) {
	out := new(Task_Prepare_Response)
	err := c.cc.Invoke(ctx, "/zedpm.plugin.TaskExecution/Prepare", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskExecutionClient) Cancel(ctx context.Context, in *Task_Cancel_Request, opts ...grpc.CallOption) (*Task_Cancel_Response, error) {
	out := new(Task_Cancel_Response)
	err := c.cc.Invoke(ctx, "/zedpm.plugin.TaskExecution/Cancel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskExecutionClient) Complete(ctx context.Context, in *Task_Complete_Request, opts ...grpc.CallOption) (*Task_Complete_Response, error) {
	out := new(Task_Complete_Response)
	err := c.cc.Invoke(ctx, "/zedpm.plugin.TaskExecution/Complete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskExecutionClient) ExecuteCheck(ctx context.Context, in *Task_Operation_Request, opts ...grpc.CallOption) (*Task_Operation_Response, error) {
	out := new(Task_Operation_Response)
	err := c.cc.Invoke(ctx, "/zedpm.plugin.TaskExecution/ExecuteCheck", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskExecutionClient) PrepareBegin(ctx context.Context, in *Task_Ref, opts ...grpc.CallOption) (*Task_SubStage_Response, error) {
	out := new(Task_SubStage_Response)
	err := c.cc.Invoke(ctx, "/zedpm.plugin.TaskExecution/PrepareBegin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskExecutionClient) ExecuteBegin(ctx context.Context, in *Task_SubStage_Request, opts ...grpc.CallOption) (*Task_Operation_Response, error) {
	out := new(Task_Operation_Response)
	err := c.cc.Invoke(ctx, "/zedpm.plugin.TaskExecution/ExecuteBegin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskExecutionClient) PrepareRun(ctx context.Context, in *Task_Ref, opts ...grpc.CallOption) (*Task_SubStage_Response, error) {
	out := new(Task_SubStage_Response)
	err := c.cc.Invoke(ctx, "/zedpm.plugin.TaskExecution/PrepareRun", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskExecutionClient) ExecuteRun(ctx context.Context, in *Task_SubStage_Request, opts ...grpc.CallOption) (*Task_Operation_Response, error) {
	out := new(Task_Operation_Response)
	err := c.cc.Invoke(ctx, "/zedpm.plugin.TaskExecution/ExecuteRun", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskExecutionClient) PrepareEnd(ctx context.Context, in *Task_Ref, opts ...grpc.CallOption) (*Task_SubStage_Response, error) {
	out := new(Task_SubStage_Response)
	err := c.cc.Invoke(ctx, "/zedpm.plugin.TaskExecution/PrepareEnd", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskExecutionClient) ExecuteEnd(ctx context.Context, in *Task_SubStage_Request, opts ...grpc.CallOption) (*Task_Operation_Response, error) {
	out := new(Task_Operation_Response)
	err := c.cc.Invoke(ctx, "/zedpm.plugin.TaskExecution/ExecuteEnd", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskExecutionClient) ExecuteFinish(ctx context.Context, in *Task_Operation_Request, opts ...grpc.CallOption) (*Task_Operation_Response, error) {
	out := new(Task_Operation_Response)
	err := c.cc.Invoke(ctx, "/zedpm.plugin.TaskExecution/ExecuteFinish", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TaskExecutionServer is the server API for TaskExecution service.
// All implementations must embed UnimplementedTaskExecutionServer
// for forward compatibility
type TaskExecutionServer interface {
	// Implements maps onto the plugin.Interface.Implements method.
	Implements(context.Context, *Task_Implements_Request) (*Task_Implements_Response, error)
	// Goal maps onto the plugin.Interface.Goal method.
	Goal(context.Context, *Task_Goal_Request) (*Task_Goal_Response, error)
	// Prepare maps onto the plugin.Interface.Prepare method.
	Prepare(context.Context, *Task_Prepare_Request) (*Task_Prepare_Response, error)
	// Cancel maps onto the plugin.Interface.Cancel method.
	Cancel(context.Context, *Task_Cancel_Request) (*Task_Cancel_Response, error)
	// Complete maps onto the plugin.Interface.Complete method.
	Complete(context.Context, *Task_Complete_Request) (*Task_Complete_Response, error)
	// ExecuteCheck maps onto the plugin.Task.Check method.
	ExecuteCheck(context.Context, *Task_Operation_Request) (*Task_Operation_Response, error)
	// PrepareBegin maps onto the plugin.Task.Begin method.
	PrepareBegin(context.Context, *Task_Ref) (*Task_SubStage_Response, error)
	// ExecuteBegin executes operations related to the plugin.Task.Begin method.
	ExecuteBegin(context.Context, *Task_SubStage_Request) (*Task_Operation_Response, error)
	// PrepareRun maps onto the plugin.Task.Run method.
	PrepareRun(context.Context, *Task_Ref) (*Task_SubStage_Response, error)
	// ExecuteRun executes operations related to the plugin.Task.Run method.
	ExecuteRun(context.Context, *Task_SubStage_Request) (*Task_Operation_Response, error)
	// PrepareEnd maps onto the plugin.Task.End method.
	PrepareEnd(context.Context, *Task_Ref) (*Task_SubStage_Response, error)
	// ExecuteEnd executes operations related to the plugin.Task.End method.
	ExecuteEnd(context.Context, *Task_SubStage_Request) (*Task_Operation_Response, error)
	// ExecuteFinish maps onto the plugin.Task.Finish method.
	ExecuteFinish(context.Context, *Task_Operation_Request) (*Task_Operation_Response, error)
	mustEmbedUnimplementedTaskExecutionServer()
}

// UnimplementedTaskExecutionServer must be embedded to have forward compatible implementations.
type UnimplementedTaskExecutionServer struct {
}

func (UnimplementedTaskExecutionServer) Implements(context.Context, *Task_Implements_Request) (*Task_Implements_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Implements not implemented")
}
func (UnimplementedTaskExecutionServer) Goal(context.Context, *Task_Goal_Request) (*Task_Goal_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Goal not implemented")
}
func (UnimplementedTaskExecutionServer) Prepare(context.Context, *Task_Prepare_Request) (*Task_Prepare_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Prepare not implemented")
}
func (UnimplementedTaskExecutionServer) Cancel(context.Context, *Task_Cancel_Request) (*Task_Cancel_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Cancel not implemented")
}
func (UnimplementedTaskExecutionServer) Complete(context.Context, *Task_Complete_Request) (*Task_Complete_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Complete not implemented")
}
func (UnimplementedTaskExecutionServer) ExecuteCheck(context.Context, *Task_Operation_Request) (*Task_Operation_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteCheck not implemented")
}
func (UnimplementedTaskExecutionServer) PrepareBegin(context.Context, *Task_Ref) (*Task_SubStage_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PrepareBegin not implemented")
}
func (UnimplementedTaskExecutionServer) ExecuteBegin(context.Context, *Task_SubStage_Request) (*Task_Operation_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteBegin not implemented")
}
func (UnimplementedTaskExecutionServer) PrepareRun(context.Context, *Task_Ref) (*Task_SubStage_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PrepareRun not implemented")
}
func (UnimplementedTaskExecutionServer) ExecuteRun(context.Context, *Task_SubStage_Request) (*Task_Operation_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteRun not implemented")
}
func (UnimplementedTaskExecutionServer) PrepareEnd(context.Context, *Task_Ref) (*Task_SubStage_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PrepareEnd not implemented")
}
func (UnimplementedTaskExecutionServer) ExecuteEnd(context.Context, *Task_SubStage_Request) (*Task_Operation_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteEnd not implemented")
}
func (UnimplementedTaskExecutionServer) ExecuteFinish(context.Context, *Task_Operation_Request) (*Task_Operation_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteFinish not implemented")
}
func (UnimplementedTaskExecutionServer) mustEmbedUnimplementedTaskExecutionServer() {}

// UnsafeTaskExecutionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TaskExecutionServer will
// result in compilation errors.
type UnsafeTaskExecutionServer interface {
	mustEmbedUnimplementedTaskExecutionServer()
}

func RegisterTaskExecutionServer(s grpc.ServiceRegistrar, srv TaskExecutionServer) {
	s.RegisterService(&TaskExecution_ServiceDesc, srv)
}

func _TaskExecution_Implements_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Task_Implements_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskExecutionServer).Implements(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/zedpm.plugin.TaskExecution/Implements",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskExecutionServer).Implements(ctx, req.(*Task_Implements_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskExecution_Goal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Task_Goal_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskExecutionServer).Goal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/zedpm.plugin.TaskExecution/Goal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskExecutionServer).Goal(ctx, req.(*Task_Goal_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskExecution_Prepare_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Task_Prepare_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskExecutionServer).Prepare(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/zedpm.plugin.TaskExecution/Prepare",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskExecutionServer).Prepare(ctx, req.(*Task_Prepare_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskExecution_Cancel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Task_Cancel_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskExecutionServer).Cancel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/zedpm.plugin.TaskExecution/Cancel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskExecutionServer).Cancel(ctx, req.(*Task_Cancel_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskExecution_Complete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Task_Complete_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskExecutionServer).Complete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/zedpm.plugin.TaskExecution/Complete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskExecutionServer).Complete(ctx, req.(*Task_Complete_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskExecution_ExecuteCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Task_Operation_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskExecutionServer).ExecuteCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/zedpm.plugin.TaskExecution/ExecuteCheck",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskExecutionServer).ExecuteCheck(ctx, req.(*Task_Operation_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskExecution_PrepareBegin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Task_Ref)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskExecutionServer).PrepareBegin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/zedpm.plugin.TaskExecution/PrepareBegin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskExecutionServer).PrepareBegin(ctx, req.(*Task_Ref))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskExecution_ExecuteBegin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Task_SubStage_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskExecutionServer).ExecuteBegin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/zedpm.plugin.TaskExecution/ExecuteBegin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskExecutionServer).ExecuteBegin(ctx, req.(*Task_SubStage_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskExecution_PrepareRun_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Task_Ref)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskExecutionServer).PrepareRun(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/zedpm.plugin.TaskExecution/PrepareRun",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskExecutionServer).PrepareRun(ctx, req.(*Task_Ref))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskExecution_ExecuteRun_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Task_SubStage_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskExecutionServer).ExecuteRun(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/zedpm.plugin.TaskExecution/ExecuteRun",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskExecutionServer).ExecuteRun(ctx, req.(*Task_SubStage_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskExecution_PrepareEnd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Task_Ref)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskExecutionServer).PrepareEnd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/zedpm.plugin.TaskExecution/PrepareEnd",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskExecutionServer).PrepareEnd(ctx, req.(*Task_Ref))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskExecution_ExecuteEnd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Task_SubStage_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskExecutionServer).ExecuteEnd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/zedpm.plugin.TaskExecution/ExecuteEnd",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskExecutionServer).ExecuteEnd(ctx, req.(*Task_SubStage_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskExecution_ExecuteFinish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Task_Operation_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskExecutionServer).ExecuteFinish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/zedpm.plugin.TaskExecution/ExecuteFinish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskExecutionServer).ExecuteFinish(ctx, req.(*Task_Operation_Request))
	}
	return interceptor(ctx, in, info, handler)
}

// TaskExecution_ServiceDesc is the grpc.ServiceDesc for TaskExecution service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TaskExecution_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "zedpm.plugin.TaskExecution",
	HandlerType: (*TaskExecutionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Implements",
			Handler:    _TaskExecution_Implements_Handler,
		},
		{
			MethodName: "Goal",
			Handler:    _TaskExecution_Goal_Handler,
		},
		{
			MethodName: "Prepare",
			Handler:    _TaskExecution_Prepare_Handler,
		},
		{
			MethodName: "Cancel",
			Handler:    _TaskExecution_Cancel_Handler,
		},
		{
			MethodName: "Complete",
			Handler:    _TaskExecution_Complete_Handler,
		},
		{
			MethodName: "ExecuteCheck",
			Handler:    _TaskExecution_ExecuteCheck_Handler,
		},
		{
			MethodName: "PrepareBegin",
			Handler:    _TaskExecution_PrepareBegin_Handler,
		},
		{
			MethodName: "ExecuteBegin",
			Handler:    _TaskExecution_ExecuteBegin_Handler,
		},
		{
			MethodName: "PrepareRun",
			Handler:    _TaskExecution_PrepareRun_Handler,
		},
		{
			MethodName: "ExecuteRun",
			Handler:    _TaskExecution_ExecuteRun_Handler,
		},
		{
			MethodName: "PrepareEnd",
			Handler:    _TaskExecution_PrepareEnd_Handler,
		},
		{
			MethodName: "ExecuteEnd",
			Handler:    _TaskExecution_ExecuteEnd_Handler,
		},
		{
			MethodName: "ExecuteFinish",
			Handler:    _TaskExecution_ExecuteFinish_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "task-interface.proto",
}
