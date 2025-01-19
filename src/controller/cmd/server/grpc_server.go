package server

import (
	"context"
	"fmt"
	"net"

	"com.kvs/deviceorbit/controller/cmd/k8s"
	"com.kvs/deviceorbit/controller/cmd/proto"
	"google.golang.org/grpc"
)

type grpcServer struct {
	podController *k8s.PodController

	listener net.Listener
	grpc     *grpc.Server

	proto.UnimplementedMobileDeviceServiceServer
}

func NewGrpcServer(podController *k8s.PodController) *grpcServer {
	return &grpcServer{podController: podController}
}

func (server *grpcServer) StartGrpcServer(port string) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return fmt.Errorf("failed to start grpc server on port: %s, cause: %s", port, err)
	}
	server.listener = listener

	server.grpc = grpc.NewServer()
	proto.RegisterMobileDeviceServiceServer(server.grpc, server)
	err = server.grpc.Serve(server.listener)
	if err != nil {
		return fmt.Errorf("failed to start grpc server on port: %s, cause: %s", port, err)
	}

	return nil
}

func (server *grpcServer) ShutdownGrpcServer() error {
	if server.grpc != nil {
		server.grpc.GracefulStop()
		server.grpc = nil
	}
	if server.listener != nil {
		err := server.listener.Close()
		if err != nil {
			return fmt.Errorf("failed to shutdown grpc server: %s", err)
		}
		server.listener = nil
	}

	return nil
}

func (s *grpcServer) CreateDevice(ctx context.Context, req *proto.DeviceRequest) (*proto.DeviceResponse, error) {
	switch req.Device.Platform {
	case proto.Platform_IOS:
		err := s.podController.CreatePod(req.Device.Serial, 1)
		if err != nil {
			return &proto.DeviceResponse{IsRunning: false}, err
		}
	case proto.Platform_ANDROID:
		err := s.podController.CreatePod(req.Device.Serial, 2)
		if err != nil {
			return &proto.DeviceResponse{IsRunning: false}, err
		}
	case proto.Platform_UNKNOWN:
		return &proto.DeviceResponse{IsRunning: false}, fmt.Errorf("unknown platform: %d", req.Device.Platform)
	}

	return &proto.DeviceResponse{IsRunning: true}, nil
}

func (s *grpcServer) DeleteDevice(ctx context.Context, req *proto.DeviceRequest) (*proto.DeviceResponse, error) {
	err := s.podController.DeletePod(req.Device.Serial)

	if err != nil {
		return &proto.DeviceResponse{IsRunning: true}, err
	}

	return &proto.DeviceResponse{IsRunning: false}, nil
}
