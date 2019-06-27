package dm

import (
	"fmt"

	"github.com/weaveworks/ignite/pkg/util"
)

// physicalDevice represents a physical block device, such as /dev/sda
type physicalDevice struct {
	file string
}

var _ blockDevice = &physicalDevice{}

func newPhysicalDevice(file string) *physicalDevice {
	return &physicalDevice{
		file: file,
	}
}

// Path to the device file
func (p *physicalDevice) Path() string {
	return p.file
}

// A physical device doesn't need to be activated beforehand, just check if it exists
func (p *physicalDevice) activate() error {
	if !util.FileExists(p.file) {
		return fmt.Errorf("physical device missing: %s", p.file)
	}

	return nil
}

// As long as the physical device is present it's active
func (p *physicalDevice) active() bool {
	return util.FileExists(p.file)
}
