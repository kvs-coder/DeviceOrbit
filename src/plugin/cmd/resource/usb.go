package resource

import (
	"fmt"

	"github.com/google/gousb"
)

type deviceType int

const (
	unknown deviceType = iota
	iOS
	android
)

type UsbInfo struct {
	Bus        int
	Address    int
	Serial     string
	DeviceType deviceType
}

func newUsbInfo(dev *gousb.Device) (UsbInfo, error) {
	desc := dev.Desc
	sn, err := dev.SerialNumber()
	if err != nil {
		return UsbInfo{}, fmt.Errorf("failed to create new usb info: %s", err)
	}

	return UsbInfo{
		Bus:        desc.Bus,
		Address:    desc.Address,
		Serial:     sn,
		DeviceType: detectDeviceType(desc),
	}, nil
}

func filterDevices(devices []UsbInfo, f func(UsbInfo) bool) []UsbInfo {
	var filtered []UsbInfo
	for _, device := range devices {
		if f(device) {
			filtered = append(filtered, device)
		}
	}
	return filtered
}

func openDevices() ([]UsbInfo, error) {
	ctx := gousb.NewContext()
	defer ctx.Close()

	iosDevices := 0
	androidDevices := 0

	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		dt := detectDeviceType(desc)
		switch dt {
		case iOS:
			iosDevices++
			return true
		case android:
			androidDevices++
			return true
		default:
			return false
		}
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open devices: %s", err)
	}

	var usbInfos []UsbInfo
	for _, dev := range devs {
		usbInfo, err := newUsbInfo(dev)
		if err != nil {
			continue
		}

		usbInfos = append(usbInfos, usbInfo)
	}

	defer func() {
		for _, d := range devs {
			d.Close()
		}
	}()

	return usbInfos, nil
}

func detectDeviceType(desc *gousb.DeviceDesc) deviceType {
	for _, cfg := range desc.Configs {
		for _, intf := range cfg.Interfaces {
			for _, ifSetting := range intf.AltSettings {
				if ifSetting.SubClass == 0x42 && ifSetting.Protocol == 1 {
					return android
				} else if ifSetting.SubClass == 0xfe && ifSetting.Protocol == 2 {
					return iOS
				}
			}
		}
	}
	return unknown
}

func handleDevice[T any](usbInfo UsbInfo, handler func(device *gousb.Device) (T, error)) (T, error) {
	var zero T
	ctx := gousb.NewContext()
	defer func(ctx *gousb.Context) {
		_ = ctx.Close()
	}(ctx)

	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return desc.Bus == usbInfo.Bus && desc.Address == usbInfo.Address
	})
	defer func() {
		for _, d := range devs {
			_ = d.Close()
		}
	}()

	if err != nil {
		return zero, fmt.Errorf("failed to open devices: %s", err)
	}

	for _, device := range devs {
		serial, err := device.SerialNumber()
		if err != nil {
			return zero, fmt.Errorf("failed to get serial number: %s", err)
		}
		if serial == usbInfo.Serial {
			return handler(device)
		}
	}
	return zero, fmt.Errorf("no device found with serial number %s", usbInfo.Serial)
}

func isDeviceConnected(usbInfo UsbInfo) bool {
	res, err := handleDevice(usbInfo, func(device *gousb.Device) (bool, error) {
		return true, nil
	})
	if err != nil {
		return false
	}
	return res
}
