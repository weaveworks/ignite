package dm

import (
	"fmt"
	"io/ioutil"
	"path"
	"strconv"

	"github.com/freddierice/go-losetup"
)

type loopDevice struct {
	losetup.Device
	file     string
	readOnly bool
	attached bool
}

var _ blockDevice = &loopDevice{}

func NewLoopDevice(file string, readOnly bool) *loopDevice {
	return &loopDevice{file: file, readOnly: readOnly}
}

func (ld *loopDevice) activate() error {
	// Don't activate twice
	if ld.active() {
		return nil
	}

	var err error
	if ld.Device, err = losetup.Attach(ld.file, 0, ld.readOnly); err != nil {
		return fmt.Errorf("failed to setup loop device for %q: %v", ld.file, err)
	} else {
		ld.attached = true
	}

	return nil
}

func (ld *loopDevice) active() bool {
	return ld.attached
}

func (ld *loopDevice) SizeSectors() (uint64, error) {
	data, err := ioutil.ReadFile(path.Join("/sys/class/block", path.Base(ld.Device.Path()), "size"))
	if err != nil {
		return 0, err
	}

	// Remove the trailing newline and parse to uint64
	return strconv.ParseUint(string(data[:len(data)-1]), 10, 64)
}
