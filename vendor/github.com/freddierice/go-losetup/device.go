package losetup

import (
	"fmt"
	"os"
)

// Device represents a loop device /dev/loop#
type Device struct {
	// device number (i.e. 7 if /dev/loop7)
	number uint64

	// flags with which to open the device with
	flags int
}

// open returns a file handle to /dev/loop# and returns an error if it cannot
// be opened.
func (device Device) open() (*os.File, error) {
	return os.OpenFile(device.Path(), device.flags, 0660)
}

// Path returns the path to the loopback device
func (device Device) Path() string {
	return fmt.Sprintf(DeviceFormatString, device.number)
}
