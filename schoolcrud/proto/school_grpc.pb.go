// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: school.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SchoolServiceClient is the client API for SchoolService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SchoolServiceClient interface {
	CreateSchool(ctx context.Context, in *School, opts ...grpc.CallOption) (*SchoolId, error)
	GetSchool(ctx context.Context, in *SchoolId, opts ...grpc.CallOption) (*School, error)
	UpdateSchool(ctx context.Context, in *School, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteSchool(ctx context.Context, in *SchoolId, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ListSchool(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (SchoolService_ListSchoolClient, error)
}

type schoolServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSchoolServiceClient(cc grpc.ClientConnInterface) SchoolServiceClient {
	return &schoolServiceClient{cc}
}

func (c *schoolServiceClient) CreateSchool(ctx context.Context, in *School, opts ...grpc.CallOption) (*SchoolId, error) {
	out := new(SchoolId)
	err := c.cc.Invoke(ctx, "/school.SchoolService/CreateSchool", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schoolServiceClient) GetSchool(ctx context.Context, in *SchoolId, opts ...grpc.CallOption) (*School, error) {
	out := new(School)
	err := c.cc.Invoke(ctx, "/school.SchoolService/GetSchool", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schoolServiceClient) UpdateSchool(ctx context.Context, in *School, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/school.SchoolService/UpdateSchool", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schoolServiceClient) DeleteSchool(ctx context.Context, in *SchoolId, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/school.SchoolService/DeleteSchool", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schoolServiceClient) ListSchool(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (SchoolService_ListSchoolClient, error) {
	stream, err := c.cc.NewStream(ctx, &SchoolService_ServiceDesc.Streams[0], "/school.SchoolService/ListSchool", opts...)
	if err != nil {
		return nil, err
	}
	x := &schoolServiceListSchoolClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SchoolService_ListSchoolClient interface {
	Recv() (*School, error)
	grpc.ClientStream
}

type schoolServiceListSchoolClient struct {
	grpc.ClientStream
}

func (x *schoolServiceListSchoolClient) Recv() (*School, error) {
	m := new(School)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SchoolServiceServer is the server API for SchoolService service.
// All implementations must embed UnimplementedSchoolServiceServer
// for forward compatibility
type SchoolServiceServer interface {
	CreateSchool(context.Context, *School) (*SchoolId, error)
	GetSchool(context.Context, *SchoolId) (*School, error)
	UpdateSchool(context.Context, *School) (*emptypb.Empty, error)
	DeleteSchool(context.Context, *SchoolId) (*emptypb.Empty, error)
	ListSchool(*emptypb.Empty, SchoolService_ListSchoolServer) error
	mustEmbedUnimplementedSchoolServiceServer()
}

// UnimplementedSchoolServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSchoolServiceServer struct {
}

func (UnimplementedSchoolServiceServer) CreateSchool(context.Context, *School) (*SchoolId, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSchool not implemented")
}
func (UnimplementedSchoolServiceServer) GetSchool(context.Context, *SchoolId) (*School, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSchool not implemented")
}
func (UnimplementedSchoolServiceServer) UpdateSchool(context.Context, *School) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSchool not implemented")
}
func (UnimplementedSchoolServiceServer) DeleteSchool(context.Context, *SchoolId) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSchool not implemented")
}
func (UnimplementedSchoolServiceServer) ListSchool(*emptypb.Empty, SchoolService_ListSchoolServer) error {
	return status.Errorf(codes.Unimplemented, "method ListSchool not implemented")
}
func (UnimplementedSchoolServiceServer) mustEmbedUnimplementedSchoolServiceServer() {}

// UnsafeSchoolServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SchoolServiceServer will
// result in compilation errors.
type UnsafeSchoolServiceServer interface {
	mustEmbedUnimplementedSchoolServiceServer()
}

func RegisterSchoolServiceServer(s grpc.ServiceRegistrar, srv SchoolServiceServer) {
	s.RegisterService(&SchoolService_ServiceDesc, srv)
}

func _SchoolService_CreateSchool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(School)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchoolServiceServer).CreateSchool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/school.SchoolService/CreateSchool",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchoolServiceServer).CreateSchool(ctx, req.(*School))
	}
	return interceptor(ctx, in, info, handler)
}

func _SchoolService_GetSchool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SchoolId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchoolServiceServer).GetSchool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/school.SchoolService/GetSchool",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchoolServiceServer).GetSchool(ctx, req.(*SchoolId))
	}
	return interceptor(ctx, in, info, handler)
}

func _SchoolService_UpdateSchool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(School)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchoolServiceServer).UpdateSchool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/school.SchoolService/UpdateSchool",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchoolServiceServer).UpdateSchool(ctx, req.(*School))
	}
	return interceptor(ctx, in, info, handler)
}

func _SchoolService_DeleteSchool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SchoolId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchoolServiceServer).DeleteSchool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/school.SchoolService/DeleteSchool",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchoolServiceServer).DeleteSchool(ctx, req.(*SchoolId))
	}
	return interceptor(ctx, in, info, handler)
}

func _SchoolService_ListSchool_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SchoolServiceServer).ListSchool(m, &schoolServiceListSchoolServer{stream})
}

type SchoolService_ListSchoolServer interface {
	Send(*School) error
	grpc.ServerStream
}

type schoolServiceListSchoolServer struct {
	grpc.ServerStream
}

func (x *schoolServiceListSchoolServer) Send(m *School) error {
	return x.ServerStream.SendMsg(m)
}

// SchoolService_ServiceDesc is the grpc.ServiceDesc for SchoolService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SchoolService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "school.SchoolService",
	HandlerType: (*SchoolServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateSchool",
			Handler:    _SchoolService_CreateSchool_Handler,
		},
		{
			MethodName: "GetSchool",
			Handler:    _SchoolService_GetSchool_Handler,
		},
		{
			MethodName: "UpdateSchool",
			Handler:    _SchoolService_UpdateSchool_Handler,
		},
		{
			MethodName: "DeleteSchool",
			Handler:    _SchoolService_DeleteSchool_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListSchool",
			Handler:       _SchoolService_ListSchool_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "school.proto",
}