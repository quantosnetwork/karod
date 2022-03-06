// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: api/proto/authorization.proto

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

// AuthorizationServerClient is the client API for AuthorizationServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthorizationServerClient interface {
	Authorize(ctx context.Context, in *AuthorizationRequest, opts ...grpc.CallOption) (*AuthorizationResponse, error)
}

type authorizationServerClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthorizationServerClient(cc grpc.ClientConnInterface) AuthorizationServerClient {
	return &authorizationServerClient{cc}
}

func (c *authorizationServerClient) Authorize(ctx context.Context, in *AuthorizationRequest, opts ...grpc.CallOption) (*AuthorizationResponse, error) {
	out := new(AuthorizationResponse)
	err := c.cc.Invoke(ctx, "/apipb.AuthorizationServer/Authorize", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthorizationServerServer is the server API for AuthorizationServer service.
// All implementations should embed UnimplementedAuthorizationServerServer
// for forward compatibility
type AuthorizationServerServer interface {
	Authorize(context.Context, *AuthorizationRequest) (*AuthorizationResponse, error)
}

// UnimplementedAuthorizationServerServer should be embedded to have forward compatible implementations.
type UnimplementedAuthorizationServerServer struct {
}

func (UnimplementedAuthorizationServerServer) Authorize(context.Context, *AuthorizationRequest) (*AuthorizationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authorize not implemented")
}

// UnsafeAuthorizationServerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthorizationServerServer will
// result in compilation errors.
type UnsafeAuthorizationServerServer interface {
	mustEmbedUnimplementedAuthorizationServerServer()
}

func RegisterAuthorizationServerServer(s grpc.ServiceRegistrar, srv AuthorizationServerServer) {
	s.RegisterService(&AuthorizationServer_ServiceDesc, srv)
}

func _AuthorizationServer_Authorize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthorizationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorizationServerServer).Authorize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/apipb.AuthorizationServer/Authorize",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorizationServerServer).Authorize(ctx, req.(*AuthorizationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthorizationServer_ServiceDesc is the grpc.ServiceDesc for AuthorizationServer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthorizationServer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "apipb.AuthorizationServer",
	HandlerType: (*AuthorizationServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Authorize",
			Handler:    _AuthorizationServer_Authorize_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/authorization.proto",
}