// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// KVServiceClient is the client API for KVService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KVServiceClient interface {
	Op(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type kVServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewKVServiceClient(cc grpc.ClientConnInterface) KVServiceClient {
	return &kVServiceClient{cc}
}

func (c *kVServiceClient) Op(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/pb.KVService/Op", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KVServiceServer is the server API for KVService service.
// All implementations must embed UnimplementedKVServiceServer
// for forward compatibility
type KVServiceServer interface {
	Op(context.Context, *Request) (*Response, error)
	mustEmbedUnimplementedKVServiceServer()
}

// UnimplementedKVServiceServer must be embedded to have forward compatible implementations.
type UnimplementedKVServiceServer struct {
}

func (UnimplementedKVServiceServer) Op(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Op not implemented")
}
func (UnimplementedKVServiceServer) mustEmbedUnimplementedKVServiceServer() {}

// UnsafeKVServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KVServiceServer will
// result in compilation errors.
type UnsafeKVServiceServer interface {
	mustEmbedUnimplementedKVServiceServer()
}

func RegisterKVServiceServer(s *grpc.Server, srv KVServiceServer) {
	s.RegisterService(&_KVService_serviceDesc, srv)
}

func _KVService_Op_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVServiceServer).Op(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.KVService/Op",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVServiceServer).Op(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _KVService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.KVService",
	HandlerType: (*KVServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Op",
			Handler:    _KVService_Op_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kv.proto",
}
