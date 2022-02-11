// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package valveapiclient

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

// ValveMatchApiServiceClient is the client API for ValveMatchApiService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ValveMatchApiServiceClient interface {
	GetNextShareCode(ctx context.Context, in *ShareCodeRequest, opts ...grpc.CallOption) (*ShareCode, error)
}

type valveMatchApiServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewValveMatchApiServiceClient(cc grpc.ClientConnInterface) ValveMatchApiServiceClient {
	return &valveMatchApiServiceClient{cc}
}

func (c *valveMatchApiServiceClient) GetNextShareCode(ctx context.Context, in *ShareCodeRequest, opts ...grpc.CallOption) (*ShareCode, error) {
	out := new(ShareCode)
	err := c.cc.Invoke(ctx, "/valveapiclient.ValveMatchApiService/GetNextShareCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ValveMatchApiServiceServer is the server API for ValveMatchApiService service.
// All implementations should embed UnimplementedValveMatchApiServiceServer
// for forward compatibility
type ValveMatchApiServiceServer interface {
	GetNextShareCode(context.Context, *ShareCodeRequest) (*ShareCode, error)
}

// UnimplementedValveMatchApiServiceServer should be embedded to have forward compatible implementations.
type UnimplementedValveMatchApiServiceServer struct {
}

func (UnimplementedValveMatchApiServiceServer) GetNextShareCode(context.Context, *ShareCodeRequest) (*ShareCode, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNextShareCode not implemented")
}

// UnsafeValveMatchApiServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ValveMatchApiServiceServer will
// result in compilation errors.
type UnsafeValveMatchApiServiceServer interface {
	mustEmbedUnimplementedValveMatchApiServiceServer()
}

func RegisterValveMatchApiServiceServer(s grpc.ServiceRegistrar, srv ValveMatchApiServiceServer) {
	s.RegisterService(&ValveMatchApiService_ServiceDesc, srv)
}

func _ValveMatchApiService_GetNextShareCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShareCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ValveMatchApiServiceServer).GetNextShareCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/valveapiclient.ValveMatchApiService/GetNextShareCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ValveMatchApiServiceServer).GetNextShareCode(ctx, req.(*ShareCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ValveMatchApiService_ServiceDesc is the grpc.ServiceDesc for ValveMatchApiService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ValveMatchApiService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "valveapiclient.ValveMatchApiService",
	HandlerType: (*ValveMatchApiServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetNextShareCode",
			Handler:    _ValveMatchApiService_GetNextShareCode_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/valveapi.proto",
}