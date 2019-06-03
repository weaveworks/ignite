package losetup

import "fmt"

// DeviceFormatString holds the format of loopback devices
const DeviceFormatString = "/dev/loop%d"

// String implements the Stringer interface for Device
func (device Device) String() string {
	return device.Path()
}

// String implements the Stringer interface for Info
func (info Info) String() string {
	return fmt.Sprintf("")
}
