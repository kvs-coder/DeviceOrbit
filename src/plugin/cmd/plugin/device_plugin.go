package plugin

import (
	"fmt"
	"net"
	"os"
	"path"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type DevicePlugin struct {
	devs         []*pluginapi.Device
	socket       string
	deviceFile   string
	resourceName string
	stop         chan interface{}
	server       *grpc.Server
}

func NewDevicePlugin(nDevices uint, deviceFilename string, resourceIdentification string, serverSock string) *DevicePlugin {
	var devs []*pluginapi.Device
	for i := uint(0); i < nDevices; i++ {
		devs = append(devs, &pluginapi.Device{
			ID:     fmt.Sprint(i),
			Health: pluginapi.Healthy,
		})
	}

	return &DevicePlugin{
		devs:         devs,
		socket:       serverSock,
		deviceFile:   deviceFilename,
		resourceName: resourceIdentification,
		stop:         make(chan interface{}),
	}
}

func dial(unixSocketPath string, timeout time.Duration) (*grpc.ClientConn, error) {
	c, err := grpc.Dial(unixSocketPath, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithTimeout(timeout),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial on unix socket %s, cause: %s", unixSocketPath, err)
	}

	return c, nil
}

func (plugin *DevicePlugin) Start(ctx context.Context) error {
	err := plugin.cleanup()
	if err != nil {
		return fmt.Errorf("failed to start plugin: %s", err)
	}

	sock, err := net.Listen("unix", plugin.socket)
	if err != nil {
		return fmt.Errorf("failed to start plugin: %s", err)
	}

	plugin.server = grpc.NewServer([]grpc.ServerOption{}...)
	pluginapi.RegisterDevicePluginServer(plugin.server, plugin)

	go plugin.server.Serve(sock)

	// Wait for server to start by launching a blocking connexion
	conn, err := dial(plugin.socket, 60*time.Second)
	if err != nil {
		return fmt.Errorf("failed to start plugin: %s", err)
	}
	conn.Close()

	return nil
}

func (plugin *DevicePlugin) Stop() error {
	if plugin.server == nil {
		return fmt.Errorf("stop failed because server is nil: %s", plugin.deviceFile)
	}

	plugin.server.Stop()
	plugin.server = nil

	close(plugin.stop)

	return plugin.cleanup()
}

func (plugin *DevicePlugin) Register(ctx context.Context, kubeletEndpoint, resourceName string) error {
	conn, err := dial(kubeletEndpoint, 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to register kubelet endpont: %s with resource name: %s, cause:%s", kubeletEndpoint, resourceName, err)
	}
	defer conn.Close()

	client := pluginapi.NewRegistrationClient(conn)
	reqt := &pluginapi.RegisterRequest{
		Version:      pluginapi.Version,
		Endpoint:     path.Base(plugin.socket),
		ResourceName: resourceName,
	}

	_, err = client.Register(ctx, reqt)
	if err != nil {
		return fmt.Errorf("failed to register kubelet endpont: %s with resource name: %s, cause:%s", kubeletEndpoint, resourceName, err)
	}
	return nil
}

func (plugin *DevicePlugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	s.Send(&pluginapi.ListAndWatchResponse{Devices: plugin.devs})

	for {
		select {
		case <-plugin.stop:
			return nil
		}
	}
}

func (plugin *DevicePlugin) Allocate(ctx context.Context, reqs *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	devs := plugin.devs
	responses := pluginapi.AllocateResponse{}
	for _, req := range reqs.ContainerRequests {
		for _, id := range req.DevicesIDs {
			if !deviceExists(devs, id) {
				return nil, fmt.Errorf("invalid allocation request: unknown device: %s", id)
			}
		}

		response := pluginapi.ContainerAllocateResponse{
			Devices: []*pluginapi.DeviceSpec{
				{
					ContainerPath: plugin.deviceFile,
					HostPath:      plugin.deviceFile,
					Permissions:   "rw",
				},
			},
		}
		responses.ContainerResponses = append(responses.ContainerResponses, &response)
	}

	return &responses, nil
}

func deviceExists(devs []*pluginapi.Device, id string) bool {
	for _, d := range devs {
		if d.ID == id {
			return true
		}
	}
	return false
}

func (plugin *DevicePlugin) cleanup() error {
	if err := os.Remove(plugin.socket); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to cleanup: %s", err)
	}

	return nil
}

func (plugin *DevicePlugin) Serve(ctx context.Context) error {
	err := plugin.Start(ctx)
	if err != nil {
		return fmt.Errorf("could not serve: %s", err)
	}

	if err := plugin.Register(ctx, pluginapi.KubeletSocket, plugin.resourceName); err != nil {
		if stopErr := plugin.Stop(); stopErr != nil {
			return fmt.Errorf("could not serve: %s; failed to stop after registration error: %s", err, stopErr)
		}
		return fmt.Errorf("could not serve: %s", err)
	}

	return nil
}

func (plugin *DevicePlugin) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{}, nil
}

func (plugin *DevicePlugin) PreStartContainer(context.Context, *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return &pluginapi.PreStartContainerResponse{}, nil
}

func (plugin *DevicePlugin) GetPreferredAllocation(context.Context, *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	return &pluginapi.PreferredAllocationResponse{}, nil
}
