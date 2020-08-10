package dmlegacy

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"strconv"

	losetup "github.com/freddierice/go-losetup"
)

// loopDevice is a helper struct for handling loopback devices for devicemapper
type loopDevice struct {
	losetup.Device
}

func newLoopDev(file string, readOnly bool) (*loopDevice, error) {
	dev, err := losetup.Attach(file, 0, readOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to setup loop device for %q: %v", file, err)
	}

	return &loopDevice{dev}, nil
}

func (ld *loopDevice) Size512K() (uint64, error) {
	data, err := ioutil.ReadFile(path.Join("/sys/class/block", path.Base(ld.Device.Path()), "size"))
	if err != nil {
		return 0, err
	}

	// Remove the trailing newline and parse to uint64
	return strconv.ParseUint(string(data[:len(data)-1]), 10, 64)
}

// dmsetup uses stdin to read multiline tables, this is a helper function for that
func runDMSetup(name string, table []byte) error {
	cmd := exec.Command(
		"dmsetup", "create",
		"--verifyudev", // if udevd is not running, dmsetup will manage the device node in /dev/mapper
		name,
	)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if _, err := stdin.Write(table); err != nil {
		return err
	}

	if err := stdin.Close(); err != nil {
		return err
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command %q exited with %q: %w", cmd.Args, out, err)
	}

	return nil
}
