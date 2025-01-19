package resource

import (
	"context"
	"fmt"

	"com.kvs/deviceorbit/plugin/cmd/plugin"
	"com.kvs/deviceorbit/plugin/cmd/utils"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type deviceResource struct {
	devicePlugin *plugin.DevicePlugin

	deviceName  string
	socketName  string
	deviceFile  string
	deviceCount uint
	usbInfo     UsbInfo
}

func newDeviceResource(usbInfo UsbInfo, requestCount uint) deviceResource {
	deviceSafeName := utils.SanitizeName(usbInfo.Serial)

	var newDevice deviceResource
	newDevice.deviceName = fmt.Sprintf("com.kvs.deviceorbit/%s", deviceSafeName)
	newDevice.socketName = fmt.Sprintf("%smdp-%s.sock", pluginapi.DevicePluginPath, deviceSafeName)
	newDevice.deviceFile = fmt.Sprintf("/dev/bus/usb/%03d/%03d", usbInfo.Bus, usbInfo.Address)
	newDevice.deviceCount = requestCount
	newDevice.usbInfo = usbInfo

	return newDevice
}

func ListDevices(requestCount uint, f func(UsbInfo) bool) ([]deviceResource, error) {
	usbInfos, err := openDevices()
	filtered := filterDevices(usbInfos, f)

	if err != nil {
		return nil, fmt.Errorf("opening devices failed: %s", err.Error())
	}

	var listDevicesAvailable []deviceResource

	for _, usbInfo := range filtered {
		newDevice := newDeviceResource(usbInfo, requestCount)
		listDevicesAvailable = append(listDevicesAvailable, newDevice)
	}

	return listDevicesAvailable, nil
}

func (res *deviceResource) Serve(ctx context.Context) error {
	if res.devicePlugin != nil {
		return fmt.Errorf("abort run, cause: plugin is already served for device: %s", res.usbInfo.Serial)
	}

	res.devicePlugin = plugin.NewDevicePlugin(res.deviceCount, res.deviceFile, res.deviceName, res.socketName)
	if err := res.devicePlugin.Serve(ctx); err != nil {
		return fmt.Errorf("failed to run device: %s, cause: %s", res.usbInfo.Serial, err)
	}

	return nil
}

func (res *deviceResource) Stop() error {
	if res.devicePlugin == nil {
		return fmt.Errorf("abort stop, cause: no plugin served for device: %s", res.usbInfo.Serial)
	}

	err := res.devicePlugin.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop device: %s, cause: %s", res.usbInfo.Serial, err)
	}
	res.devicePlugin = nil

	return nil
}

func (res deviceResource) SanitizedUsbSerial() string {
	return utils.SanitizeName(res.usbInfo.Serial)
}

func (res deviceResource) IsConnected() bool {
	return isDeviceConnected(res.usbInfo)
}

func (res deviceResource) Exists(path string) bool {
	exists := utils.FileExists(res.deviceFile)
	return path == res.deviceFile && !exists
}

func (res deviceResource) String() string {
	return fmt.Sprintf("DeviceResource(Serial: %s, File: %s, Socket: %s)", res.SanitizedUsbSerial(), res.deviceFile, res.socketName)
}
