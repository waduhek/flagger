// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: proto/flagpb/flag.proto

package flagpb

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
	Flag_CreateFlag_FullMethodName       = "/flagpb.Flag/CreateFlag"
	Flag_UpdateFlagStatus_FullMethodName = "/flagpb.Flag/UpdateFlagStatus"
)

// FlagClient is the client API for Flag service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FlagClient interface {
	// CreateFlag creates a new flag under the project. The flag will be created
	// in all environments and will be active by default.
	CreateFlag(ctx context.Context, in *CreateFlagRequest, opts ...grpc.CallOption) (*CreateFlagResponse, error)
	// UpdateFlagStatus updates the status of the flag in the provided
	// environment.
	UpdateFlagStatus(ctx context.Context, in *UpdateFlagStatusRequest, opts ...grpc.CallOption) (*UpdateFlagStatusResponse, error)
}

type flagClient struct {
	cc grpc.ClientConnInterface
}

func NewFlagClient(cc grpc.ClientConnInterface) FlagClient {
	return &flagClient{cc}
}

func (c *flagClient) CreateFlag(ctx context.Context, in *CreateFlagRequest, opts ...grpc.CallOption) (*CreateFlagResponse, error) {
	out := new(CreateFlagResponse)
	err := c.cc.Invoke(ctx, Flag_CreateFlag_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *flagClient) UpdateFlagStatus(ctx context.Context, in *UpdateFlagStatusRequest, opts ...grpc.CallOption) (*UpdateFlagStatusResponse, error) {
	out := new(UpdateFlagStatusResponse)
	err := c.cc.Invoke(ctx, Flag_UpdateFlagStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FlagServer is the server API for Flag service.
// All implementations must embed UnimplementedFlagServer
// for forward compatibility
type FlagServer interface {
	// CreateFlag creates a new flag under the project. The flag will be created
	// in all environments and will be active by default.
	CreateFlag(context.Context, *CreateFlagRequest) (*CreateFlagResponse, error)
	// UpdateFlagStatus updates the status of the flag in the provided
	// environment.
	UpdateFlagStatus(context.Context, *UpdateFlagStatusRequest) (*UpdateFlagStatusResponse, error)
	mustEmbedUnimplementedFlagServer()
}

// UnimplementedFlagServer must be embedded to have forward compatible implementations.
type UnimplementedFlagServer struct {
}

func (UnimplementedFlagServer) CreateFlag(context.Context, *CreateFlagRequest) (*CreateFlagResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFlag not implemented")
}
func (UnimplementedFlagServer) UpdateFlagStatus(context.Context, *UpdateFlagStatusRequest) (*UpdateFlagStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateFlagStatus not implemented")
}
func (UnimplementedFlagServer) mustEmbedUnimplementedFlagServer() {}

// UnsafeFlagServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FlagServer will
// result in compilation errors.
type UnsafeFlagServer interface {
	mustEmbedUnimplementedFlagServer()
}

func RegisterFlagServer(s grpc.ServiceRegistrar, srv FlagServer) {
	s.RegisterService(&Flag_ServiceDesc, srv)
}

func _Flag_CreateFlag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateFlagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FlagServer).CreateFlag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Flag_CreateFlag_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FlagServer).CreateFlag(ctx, req.(*CreateFlagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Flag_UpdateFlagStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateFlagStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FlagServer).UpdateFlagStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Flag_UpdateFlagStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FlagServer).UpdateFlagStatus(ctx, req.(*UpdateFlagStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Flag_ServiceDesc is the grpc.ServiceDesc for Flag service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Flag_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "flagpb.Flag",
	HandlerType: (*FlagServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateFlag",
			Handler:    _Flag_CreateFlag_Handler,
		},
		{
			MethodName: "UpdateFlagStatus",
			Handler:    _Flag_UpdateFlagStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/flagpb/flag.proto",
}
