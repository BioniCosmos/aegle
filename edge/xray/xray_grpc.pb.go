// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: edge/xray/xray.proto

package xray

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

const (
	Xray_AddInbound_FullMethodName    = "/xray.Xray/AddInbound"
	Xray_RemoveInbound_FullMethodName = "/xray.Xray/RemoveInbound"
	Xray_AddUser_FullMethodName       = "/xray.Xray/AddUser"
	Xray_RemoveUser_FullMethodName    = "/xray.Xray/RemoveUser"
)

// XrayClient is the client API for Xray service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type XrayClient interface {
	AddInbound(ctx context.Context, in *AddInboundRequest, opts ...grpc.CallOption) (*Response, error)
	RemoveInbound(ctx context.Context, in *RemoveInboundRequest, opts ...grpc.CallOption) (*Response, error)
	AddUser(ctx context.Context, in *AddUserRequest, opts ...grpc.CallOption) (*Response, error)
	RemoveUser(ctx context.Context, in *RemoveUserRequest, opts ...grpc.CallOption) (*Response, error)
}

type xrayClient struct {
	cc grpc.ClientConnInterface
}

func NewXrayClient(cc grpc.ClientConnInterface) XrayClient {
	return &xrayClient{cc}
}

func (c *xrayClient) AddInbound(ctx context.Context, in *AddInboundRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, Xray_AddInbound_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *xrayClient) RemoveInbound(ctx context.Context, in *RemoveInboundRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, Xray_RemoveInbound_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *xrayClient) AddUser(ctx context.Context, in *AddUserRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, Xray_AddUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *xrayClient) RemoveUser(ctx context.Context, in *RemoveUserRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, Xray_RemoveUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// XrayServer is the server API for Xray service.
// All implementations must embed UnimplementedXrayServer
// for forward compatibility
type XrayServer interface {
	AddInbound(context.Context, *AddInboundRequest) (*Response, error)
	RemoveInbound(context.Context, *RemoveInboundRequest) (*Response, error)
	AddUser(context.Context, *AddUserRequest) (*Response, error)
	RemoveUser(context.Context, *RemoveUserRequest) (*Response, error)
	mustEmbedUnimplementedXrayServer()
}

// UnimplementedXrayServer must be embedded to have forward compatible implementations.
type UnimplementedXrayServer struct {
}

func (UnimplementedXrayServer) AddInbound(context.Context, *AddInboundRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddInbound not implemented")
}
func (UnimplementedXrayServer) RemoveInbound(context.Context, *RemoveInboundRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveInbound not implemented")
}
func (UnimplementedXrayServer) AddUser(context.Context, *AddUserRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUser not implemented")
}
func (UnimplementedXrayServer) RemoveUser(context.Context, *RemoveUserRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveUser not implemented")
}
func (UnimplementedXrayServer) mustEmbedUnimplementedXrayServer() {}

// UnsafeXrayServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to XrayServer will
// result in compilation errors.
type UnsafeXrayServer interface {
	mustEmbedUnimplementedXrayServer()
}

func RegisterXrayServer(s grpc.ServiceRegistrar, srv XrayServer) {
	s.RegisterService(&Xray_ServiceDesc, srv)
}

func _Xray_AddInbound_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddInboundRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(XrayServer).AddInbound(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Xray_AddInbound_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(XrayServer).AddInbound(ctx, req.(*AddInboundRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Xray_RemoveInbound_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveInboundRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(XrayServer).RemoveInbound(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Xray_RemoveInbound_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(XrayServer).RemoveInbound(ctx, req.(*RemoveInboundRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Xray_AddUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(XrayServer).AddUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Xray_AddUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(XrayServer).AddUser(ctx, req.(*AddUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Xray_RemoveUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(XrayServer).RemoveUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Xray_RemoveUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(XrayServer).RemoveUser(ctx, req.(*RemoveUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Xray_ServiceDesc is the grpc.ServiceDesc for Xray service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Xray_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "xray.Xray",
	HandlerType: (*XrayServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddInbound",
			Handler:    _Xray_AddInbound_Handler,
		},
		{
			MethodName: "RemoveInbound",
			Handler:    _Xray_RemoveInbound_Handler,
		},
		{
			MethodName: "AddUser",
			Handler:    _Xray_AddUser_Handler,
		},
		{
			MethodName: "RemoveUser",
			Handler:    _Xray_RemoveUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "edge/xray/xray.proto",
}
