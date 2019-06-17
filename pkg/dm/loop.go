package dm

import (
	"fmt"
	"github.com/freddierice/go-losetup"
	"io/ioutil"
	"path"
	"strconv"
)

type loopDevice struct {
	losetup.Device
	file string
}

var _ blockDevice = &loopDevice{}

func NewLoopDevice(file string) *loopDevice {
	return &loopDevice{file: file}
}

func (ld *loopDevice) Attach(readOnly bool) error {
	var err error
	if ld.Device, err = losetup.Attach(ld.file, 0, readOnly); err != nil {
		return fmt.Errorf("failed to setup loop device for %q: %v", ld.file, err)
	}

	return nil
}

func (ld *loopDevice) Size512K() (uint64, error) {
	data, err := ioutil.ReadFile(path.Join("/sys/class/block", path.Base(ld.Device.Path()), "size"))
	if err != nil {
		return 0, err
	}

	// Remove the trailing newline and parse to uint64
	return strconv.ParseUint(string(data[:len(data)-1]), 10, 64)
}
