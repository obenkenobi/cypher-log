// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.4
// source: userkeypb/userkey.proto

package userkeypb

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

// UserKeyServiceClient is the client API for UserKeyService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserKeyServiceClient interface {
	GetKeyFromSession(ctx context.Context, in *UserKeySession, opts ...grpc.CallOption) (*UserKey, error)
}

type userKeyServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserKeyServiceClient(cc grpc.ClientConnInterface) UserKeyServiceClient {
	return &userKeyServiceClient{cc}
}

func (c *userKeyServiceClient) GetKeyFromSession(ctx context.Context, in *UserKeySession, opts ...grpc.CallOption) (*UserKey, error) {
	out := new(UserKey)
	err := c.cc.Invoke(ctx, "/UserKeyService/GetKeyFromSession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserKeyServiceServer is the server API for UserKeyService service.
// All implementations must embed UnimplementedUserKeyServiceServer
// for forward compatibility
type UserKeyServiceServer interface {
	GetKeyFromSession(context.Context, *UserKeySession) (*UserKey, error)
	mustEmbedUnimplementedUserKeyServiceServer()
}

// UnimplementedUserKeyServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserKeyServiceServer struct {
}

func (UnimplementedUserKeyServiceServer) GetKeyFromSession(context.Context, *UserKeySession) (*UserKey, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKeyFromSession not implemented")
}
func (UnimplementedUserKeyServiceServer) mustEmbedUnimplementedUserKeyServiceServer() {}

// UnsafeUserKeyServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserKeyServiceServer will
// result in compilation errors.
type UnsafeUserKeyServiceServer interface {
	mustEmbedUnimplementedUserKeyServiceServer()
}

func RegisterUserKeyServiceServer(s grpc.ServiceRegistrar, srv UserKeyServiceServer) {
	s.RegisterService(&UserKeyService_ServiceDesc, srv)
}

func _UserKeyService_GetKeyFromSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserKeySession)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserKeyServiceServer).GetKeyFromSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/UserKeyService/GetKeyFromSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserKeyServiceServer).GetKeyFromSession(ctx, req.(*UserKeySession))
	}
	return interceptor(ctx, in, info, handler)
}

// UserKeyService_ServiceDesc is the grpc.ServiceDesc for UserKeyService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserKeyService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "UserKeyService",
	HandlerType: (*UserKeyServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetKeyFromSession",
			Handler:    _UserKeyService_GetKeyFromSession_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "userkeypb/userkey.proto",
}
