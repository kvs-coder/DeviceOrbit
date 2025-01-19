package client

import (
	"context"
	"fmt"
	"time"

	"com.kvs/deviceorbit/plugin/cmd/proto"
	"google.golang.org/grpc"
)

type GrpcClient struct {
	conn   *grpc.ClientConn
	ctx    context.Context
	client proto.MobileDeviceServiceClient
}

func NewGrpcCLient(ctx context.Context) (*GrpcClient, error) {
	serverAddress := "mobile-device-controller.default.svc.cluster.local:50051"
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("grpc client did not connect: %s", err)
	}

	client := proto.NewMobileDeviceServiceClient(conn)

	return &GrpcClient{conn, ctx, client}, nil
}

func (gc *GrpcClient) CreateDevice(deviceSerial string, platform int) (bool, error) {
	req := &proto.DeviceRequest{
		Device: &proto.Device{
			Serial:   deviceSerial,
			Platform: proto.Platform(platform),
		},
	}

	ctx, cancel := context.WithTimeout(gc.ctx, time.Second*5)
	defer cancel()

	res, err := gc.client.CreateDevice(ctx, req)
	if err != nil {
		return false, fmt.Errorf("failed to create k8s device: %s", err)
	}

	return res.IsRunning, nil
}

func (gc *GrpcClient) DeleteDevice(deviceSerial string) (bool, error) {
	req := &proto.DeviceRequest{
		Device: &proto.Device{
			Serial: deviceSerial,
		},
	}

	ctx, cancel := context.WithTimeout(gc.ctx, time.Second*5)
	defer cancel()

	res, err := gc.client.DeleteDevice(ctx, req)
	if err != nil {
		return false, fmt.Errorf("failed to delete k8s device: %s", err)
	}

	return res.IsRunning, nil
}

func (gc *GrpcClient) Close() error {
	return gc.conn.Close()
}
