package resource

import (
	"context"
	"fmt"
	"slices"

	"com.kvs/deviceorbit/plugin/cmd/client"
	"com.kvs/deviceorbit/plugin/cmd/utils"
	"github.com/radovskyb/watcher"
)

type DeviceRegistry struct {
	client *client.GrpcClient

	runningDevices []deviceResource
	requestCount   uint
}

func NewDeviceRegistry(client *client.GrpcClient, requestCount uint) *DeviceRegistry {
	return &DeviceRegistry{client: client, requestCount: requestCount}
}

func (d *DeviceRegistry) HandleDeviceEvent(ctx context.Context, event watcher.Event) (deviceResource, error) {
	if !utils.HasMatch(event.Path) {
		return deviceResource{}, nil
	}

	switch event.Op {
	case watcher.Create:
		device, err := d.findPluggedDevice(event.Path)
		if err != nil {
			return deviceResource{}, fmt.Errorf("failed to handle event 'Create' on path: %s, cause: %s", event.Path, err)
		}

		err = d.runDevice(ctx, device)
		if err != nil {
			return deviceResource{}, fmt.Errorf("failed to handle event 'Create' on path: %s, cause: %s", event.Path, err)
		}

		return device, nil

	case watcher.Remove:
		device, err := d.findUnpluggedDevice(event.Path)
		if err != nil {
			return deviceResource{}, fmt.Errorf("failed to handle event 'Remove' on path: %s, cause: %s", event.Path, err)
		}

		err = d.stopDevice(device)
		if err != nil {
			return deviceResource{}, fmt.Errorf("failed to handle event 'Remove' on path: %s, cause: %s", event.Path, err)
		}

		return device, nil
	}

	return deviceResource{}, nil
}

func (d *DeviceRegistry) IsEmptyResource(device deviceResource) bool {
	return device == (deviceResource{})
}

func (d *DeviceRegistry) ListRunningDevices() []deviceResource {
	return d.runningDevices
}

func (d *DeviceRegistry) RunDevices(ctx context.Context) error {
	discoveredDevices, err := ListDevices(d.requestCount, func(ui UsbInfo) bool { return true })
	if err != nil {
		return fmt.Errorf("failed to run discovered devices: %s", err)
	}

	if len(discoveredDevices) == 0 {
		return fmt.Errorf("no devices plugged in")
	}

	d.runningDevices = make([]deviceResource, 0)

	for _, device := range discoveredDevices {
		err := d.runDevice(ctx, device)
		if err != nil {
			continue
		}
	}

	return nil
}

func (d *DeviceRegistry) StopDevices() error {
	if len(d.runningDevices) == 0 {
		return fmt.Errorf("no devices running")
	}

	discoveredDevices, err := ListDevices(d.requestCount, func(ui UsbInfo) bool { return true })
	if err != nil {
		return fmt.Errorf("failed to stop all devices: %s", err)
	}

	for _, device := range discoveredDevices {
		err := d.stopDevice(device)
		if err != nil {
			continue
		}
	}

	return nil
}

func (d *DeviceRegistry) stopDevice(runningDevice deviceResource) error {
	_, err := d.client.DeleteDevice(runningDevice.usbInfo.Serial)
	if err != nil {
		return fmt.Errorf("failed to stop device: %s, cause: %s", runningDevice.usbInfo.Serial, err)
	}

	err = runningDevice.Stop()
	if err != nil {
		return err
	}

	d.runningDevices = slices.DeleteFunc(d.runningDevices, func(device deviceResource) bool {
		return runningDevice.usbInfo.Serial == device.usbInfo.Serial
	})

	return nil
}

func (d *DeviceRegistry) runDevice(ctx context.Context, device deviceResource) error {
	err := device.Serve(ctx)
	if err != nil {
		return err
	}

	_, err = d.client.CreateDevice(device.SanitizedUsbSerial(), int(device.usbInfo.DeviceType))
	if err != nil {
		return fmt.Errorf("failed to run device: %s, cause: %s", device.usbInfo.Serial, err)
	}

	d.runningDevices = append(d.runningDevices, device)

	return nil
}

func (d *DeviceRegistry) findUnpluggedDevice(path string) (deviceResource, error) {
	for _, device := range d.runningDevices {
		if device.Exists(path) {
			return device, nil
		}
	}

	return deviceResource{}, fmt.Errorf("no such device in path: %s", path)
}

func (d *DeviceRegistry) findPluggedDevice(path string) (deviceResource, error) {
	usbDevice, err := utils.FindUsbDeviceMatch(path)
	if err != nil {
		return deviceResource{}, fmt.Errorf("failed to find plugged device: %s", err)
	}

	newDevices, err := ListDevices(d.requestCount, func(ui UsbInfo) bool {
		return ui.Bus == usbDevice.Bus && ui.Address == usbDevice.Address
	})
	if err != nil {
		return deviceResource{}, fmt.Errorf("failed to find plugged device: %s", err)
	}

	return newDevices[0], nil
}
