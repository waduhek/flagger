// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: proto/environmentpb/environment.proto

package environmentpb

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
	Environment_CreateEnvironment_FullMethodName = "/environmentpb.Environment/CreateEnvironment"
)

// EnvironmentClient is the client API for Environment service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EnvironmentClient interface {
	// CreateEnvironment creates a new environment for a project.
	CreateEnvironment(ctx context.Context, in *CreateEnvironmentRequest, opts ...grpc.CallOption) (*CreateEnvironmentResponse, error)
}

type environmentClient struct {
	cc grpc.ClientConnInterface
}

func NewEnvironmentClient(cc grpc.ClientConnInterface) EnvironmentClient {
	return &environmentClient{cc}
}

func (c *environmentClient) CreateEnvironment(ctx context.Context, in *CreateEnvironmentRequest, opts ...grpc.CallOption) (*CreateEnvironmentResponse, error) {
	out := new(CreateEnvironmentResponse)
	err := c.cc.Invoke(ctx, Environment_CreateEnvironment_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EnvironmentServer is the server API for Environment service.
// All implementations must embed UnimplementedEnvironmentServer
// for forward compatibility
type EnvironmentServer interface {
	// CreateEnvironment creates a new environment for a project.
	CreateEnvironment(context.Context, *CreateEnvironmentRequest) (*CreateEnvironmentResponse, error)
	mustEmbedUnimplementedEnvironmentServer()
}

// UnimplementedEnvironmentServer must be embedded to have forward compatible implementations.
type UnimplementedEnvironmentServer struct {
}

func (UnimplementedEnvironmentServer) CreateEnvironment(context.Context, *CreateEnvironmentRequest) (*CreateEnvironmentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateEnvironment not implemented")
}
func (UnimplementedEnvironmentServer) mustEmbedUnimplementedEnvironmentServer() {}

// UnsafeEnvironmentServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EnvironmentServer will
// result in compilation errors.
type UnsafeEnvironmentServer interface {
	mustEmbedUnimplementedEnvironmentServer()
}

func RegisterEnvironmentServer(s grpc.ServiceRegistrar, srv EnvironmentServer) {
	s.RegisterService(&Environment_ServiceDesc, srv)
}

func _Environment_CreateEnvironment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateEnvironmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnvironmentServer).CreateEnvironment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Environment_CreateEnvironment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnvironmentServer).CreateEnvironment(ctx, req.(*CreateEnvironmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Environment_ServiceDesc is the grpc.ServiceDesc for Environment service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Environment_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "environmentpb.Environment",
	HandlerType: (*EnvironmentServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateEnvironment",
			Handler:    _Environment_CreateEnvironment_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/environmentpb/environment.proto",
}
