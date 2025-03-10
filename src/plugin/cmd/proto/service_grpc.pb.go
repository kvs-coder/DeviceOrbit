// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: service.proto

package proto

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
	MobileDeviceService_CreateDevice_FullMethodName = "/com.kvs.mobile_device_controller.MobileDeviceService/CreateDevice"
	MobileDeviceService_DeleteDevice_FullMethodName = "/com.kvs.mobile_device_controller.MobileDeviceService/DeleteDevice"
)

// MobileDeviceServiceClient is the client API for MobileDeviceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MobileDeviceServiceClient interface {
	CreateDevice(ctx context.Context, in *DeviceRequest, opts ...grpc.CallOption) (*DeviceResponse, error)
	DeleteDevice(ctx context.Context, in *DeviceRequest, opts ...grpc.CallOption) (*DeviceResponse, error)
}

type mobileDeviceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMobileDeviceServiceClient(cc grpc.ClientConnInterface) MobileDeviceServiceClient {
	return &mobileDeviceServiceClient{cc}
}

func (c *mobileDeviceServiceClient) CreateDevice(ctx context.Context, in *DeviceRequest, opts ...grpc.CallOption) (*DeviceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeviceResponse)
	err := c.cc.Invoke(ctx, MobileDeviceService_CreateDevice_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mobileDeviceServiceClient) DeleteDevice(ctx context.Context, in *DeviceRequest, opts ...grpc.CallOption) (*DeviceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeviceResponse)
	err := c.cc.Invoke(ctx, MobileDeviceService_DeleteDevice_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MobileDeviceServiceServer is the server API for MobileDeviceService service.
// All implementations must embed UnimplementedMobileDeviceServiceServer
// for forward compatibility.
type MobileDeviceServiceServer interface {
	CreateDevice(context.Context, *DeviceRequest) (*DeviceResponse, error)
	DeleteDevice(context.Context, *DeviceRequest) (*DeviceResponse, error)
	mustEmbedUnimplementedMobileDeviceServiceServer()
}

// UnimplementedMobileDeviceServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMobileDeviceServiceServer struct{}

func (UnimplementedMobileDeviceServiceServer) CreateDevice(context.Context, *DeviceRequest) (*DeviceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateDevice not implemented")
}
func (UnimplementedMobileDeviceServiceServer) DeleteDevice(context.Context, *DeviceRequest) (*DeviceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteDevice not implemented")
}
func (UnimplementedMobileDeviceServiceServer) mustEmbedUnimplementedMobileDeviceServiceServer() {}
func (UnimplementedMobileDeviceServiceServer) testEmbeddedByValue()                             {}

// UnsafeMobileDeviceServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MobileDeviceServiceServer will
// result in compilation errors.
type UnsafeMobileDeviceServiceServer interface {
	mustEmbedUnimplementedMobileDeviceServiceServer()
}

func RegisterMobileDeviceServiceServer(s grpc.ServiceRegistrar, srv MobileDeviceServiceServer) {
	// If the following call pancis, it indicates UnimplementedMobileDeviceServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MobileDeviceService_ServiceDesc, srv)
}

func _MobileDeviceService_CreateDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MobileDeviceServiceServer).CreateDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MobileDeviceService_CreateDevice_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MobileDeviceServiceServer).CreateDevice(ctx, req.(*DeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MobileDeviceService_DeleteDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MobileDeviceServiceServer).DeleteDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MobileDeviceService_DeleteDevice_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MobileDeviceServiceServer).DeleteDevice(ctx, req.(*DeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MobileDeviceService_ServiceDesc is the grpc.ServiceDesc for MobileDeviceService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MobileDeviceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "com.kvs.mobile_device_controller.MobileDeviceService",
	HandlerType: (*MobileDeviceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateDevice",
			Handler:    _MobileDeviceService_CreateDevice_Handler,
		},
		{
			MethodName: "DeleteDevice",
			Handler:    _MobileDeviceService_DeleteDevice_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}
