// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: cli/v1/cliserver.proto

package cliserver

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Cli_AcceptOrderGrpc_FullMethodName  = "/cliserver.cli/AcceptOrderGrpc"
	Cli_AcceptReturnGrpc_FullMethodName = "/cliserver.cli/AcceptReturnGrpc"
	Cli_GiveOrderGrpc_FullMethodName    = "/cliserver.cli/GiveOrderGrpc"
	Cli_ReturnOrderGrpc_FullMethodName  = "/cliserver.cli/ReturnOrderGrpc"
	Cli_UserOrdersGrpc_FullMethodName   = "/cliserver.cli/UserOrdersGrpc"
	Cli_UserReturnsGrpc_FullMethodName  = "/cliserver.cli/UserReturnsGrpc"
)

// CliClient is the client API for Cli service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CliClient interface {
	AcceptOrderGrpc(ctx context.Context, in *AcceptOrderRequest, opts ...grpc.CallOption) (*AcceptOrderResponse, error)
	AcceptReturnGrpc(ctx context.Context, in *AcceptReturnRequest, opts ...grpc.CallOption) (*AcceptReturnResponse, error)
	GiveOrderGrpc(ctx context.Context, in *GiveOrderRequest, opts ...grpc.CallOption) (*GiveOrderResponse, error)
	ReturnOrderGrpc(ctx context.Context, in *ReturnOrderRequest, opts ...grpc.CallOption) (*ReturnOrderResponse, error)
	UserOrdersGrpc(ctx context.Context, in *UserOrdersRequest, opts ...grpc.CallOption) (*UserOrdersResponse, error)
	UserReturnsGrpc(ctx context.Context, in *UserReturnsRequest, opts ...grpc.CallOption) (*UserReturnsResponse, error)
}

type cliClient struct {
	cc grpc.ClientConnInterface
}

func NewCliClient(cc grpc.ClientConnInterface) CliClient {
	return &cliClient{cc}
}

func (c *cliClient) AcceptOrderGrpc(ctx context.Context, in *AcceptOrderRequest, opts ...grpc.CallOption) (*AcceptOrderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AcceptOrderResponse)
	err := c.cc.Invoke(ctx, Cli_AcceptOrderGrpc_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cliClient) AcceptReturnGrpc(ctx context.Context, in *AcceptReturnRequest, opts ...grpc.CallOption) (*AcceptReturnResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AcceptReturnResponse)
	err := c.cc.Invoke(ctx, Cli_AcceptReturnGrpc_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cliClient) GiveOrderGrpc(ctx context.Context, in *GiveOrderRequest, opts ...grpc.CallOption) (*GiveOrderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GiveOrderResponse)
	err := c.cc.Invoke(ctx, Cli_GiveOrderGrpc_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cliClient) ReturnOrderGrpc(ctx context.Context, in *ReturnOrderRequest, opts ...grpc.CallOption) (*ReturnOrderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReturnOrderResponse)
	err := c.cc.Invoke(ctx, Cli_ReturnOrderGrpc_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cliClient) UserOrdersGrpc(ctx context.Context, in *UserOrdersRequest, opts ...grpc.CallOption) (*UserOrdersResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserOrdersResponse)
	err := c.cc.Invoke(ctx, Cli_UserOrdersGrpc_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cliClient) UserReturnsGrpc(ctx context.Context, in *UserReturnsRequest, opts ...grpc.CallOption) (*UserReturnsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserReturnsResponse)
	err := c.cc.Invoke(ctx, Cli_UserReturnsGrpc_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CliServer is the server API for Cli service.
// All implementations must embed UnimplementedCliServer
// for forward compatibility.
type CliServer interface {
	AcceptOrderGrpc(context.Context, *AcceptOrderRequest) (*AcceptOrderResponse, error)
	AcceptReturnGrpc(context.Context, *AcceptReturnRequest) (*AcceptReturnResponse, error)
	GiveOrderGrpc(context.Context, *GiveOrderRequest) (*GiveOrderResponse, error)
	ReturnOrderGrpc(context.Context, *ReturnOrderRequest) (*ReturnOrderResponse, error)
	UserOrdersGrpc(context.Context, *UserOrdersRequest) (*UserOrdersResponse, error)
	UserReturnsGrpc(context.Context, *UserReturnsRequest) (*UserReturnsResponse, error)
	mustEmbedUnimplementedCliServer()
}

// UnimplementedCliServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCliServer struct{}

func (UnimplementedCliServer) AcceptOrderGrpc(context.Context, *AcceptOrderRequest) (*AcceptOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AcceptOrderGrpc not implemented")
}
func (UnimplementedCliServer) AcceptReturnGrpc(context.Context, *AcceptReturnRequest) (*AcceptReturnResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AcceptReturnGrpc not implemented")
}
func (UnimplementedCliServer) GiveOrderGrpc(context.Context, *GiveOrderRequest) (*GiveOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GiveOrderGrpc not implemented")
}
func (UnimplementedCliServer) ReturnOrderGrpc(context.Context, *ReturnOrderRequest) (*ReturnOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReturnOrderGrpc not implemented")
}
func (UnimplementedCliServer) UserOrdersGrpc(context.Context, *UserOrdersRequest) (*UserOrdersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserOrdersGrpc not implemented")
}
func (UnimplementedCliServer) UserReturnsGrpc(context.Context, *UserReturnsRequest) (*UserReturnsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserReturnsGrpc not implemented")
}
func (UnimplementedCliServer) mustEmbedUnimplementedCliServer() {}
func (UnimplementedCliServer) testEmbeddedByValue()             {}

// UnsafeCliServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CliServer will
// result in compilation errors.
type UnsafeCliServer interface {
	mustEmbedUnimplementedCliServer()
}

func RegisterCliServer(s grpc.ServiceRegistrar, srv CliServer) {
	// If the following call pancis, it indicates UnimplementedCliServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Cli_ServiceDesc, srv)
}

func _Cli_AcceptOrderGrpc_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AcceptOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CliServer).AcceptOrderGrpc(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cli_AcceptOrderGrpc_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CliServer).AcceptOrderGrpc(ctx, req.(*AcceptOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cli_AcceptReturnGrpc_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AcceptReturnRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CliServer).AcceptReturnGrpc(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cli_AcceptReturnGrpc_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CliServer).AcceptReturnGrpc(ctx, req.(*AcceptReturnRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cli_GiveOrderGrpc_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GiveOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CliServer).GiveOrderGrpc(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cli_GiveOrderGrpc_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CliServer).GiveOrderGrpc(ctx, req.(*GiveOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cli_ReturnOrderGrpc_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReturnOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CliServer).ReturnOrderGrpc(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cli_ReturnOrderGrpc_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CliServer).ReturnOrderGrpc(ctx, req.(*ReturnOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cli_UserOrdersGrpc_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserOrdersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CliServer).UserOrdersGrpc(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cli_UserOrdersGrpc_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CliServer).UserOrdersGrpc(ctx, req.(*UserOrdersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cli_UserReturnsGrpc_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserReturnsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CliServer).UserReturnsGrpc(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cli_UserReturnsGrpc_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CliServer).UserReturnsGrpc(ctx, req.(*UserReturnsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Cli_ServiceDesc is the grpc.ServiceDesc for Cli service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Cli_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cliserver.cli",
	HandlerType: (*CliServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AcceptOrderGrpc",
			Handler:    _Cli_AcceptOrderGrpc_Handler,
		},
		{
			MethodName: "AcceptReturnGrpc",
			Handler:    _Cli_AcceptReturnGrpc_Handler,
		},
		{
			MethodName: "GiveOrderGrpc",
			Handler:    _Cli_GiveOrderGrpc_Handler,
		},
		{
			MethodName: "ReturnOrderGrpc",
			Handler:    _Cli_ReturnOrderGrpc_Handler,
		},
		{
			MethodName: "UserOrdersGrpc",
			Handler:    _Cli_UserOrdersGrpc_Handler,
		},
		{
			MethodName: "UserReturnsGrpc",
			Handler:    _Cli_UserReturnsGrpc_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cli/v1/cliserver.proto",
}
