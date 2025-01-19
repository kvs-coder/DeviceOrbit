package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func HasMatch(path string) bool {
	re := regexp.MustCompile(`^/dev/bus/usb/\d{3}/\d{3}$`)
	return re.MatchString(path)
}

type usbDevice struct {
	Bus     int
	Address int
}

func FindUsbDeviceMatch(path string) (usbDevice, error) {
	re := regexp.MustCompile(`^/dev/bus/usb/(\d{3})/(\d{3})$`)

	matches := re.FindStringSubmatch(path)
	if len(matches) != 3 {
		return usbDevice{}, fmt.Errorf("invalid device path format: %s", path)
	}

	bus, err := strconv.Atoi(matches[1])
	if err != nil {
		return usbDevice{}, err
	}
	address, err := strconv.Atoi(matches[2])
	if err != nil {
		return usbDevice{}, err
	}

	return usbDevice{bus, address}, nil
}

func SanitizeName(path string) string {
	sanitizeChar := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return r
		case r >= 'a' && r <= 'z':
			return r
		case r >= '0' && r <= '9':
			return r
		case r == '_':
			return r
		case r == '-':
			return r
		}
		return -1
	}
	return strings.Map(sanitizeChar, path)
}
