package dm

import (
	"log"
	"os"

	"github.com/weaveworks/ignite/pkg/util"
)

type blockDevice interface {
	Path() string
	activate() error
	active() bool
}

func dmsetup(args ...string) error {
	log.Printf("Running dmsetup: %q\n", args)
	_, err := util.ExecuteCommand("dmsetup", args...)
	return err
}

func mkfs(device blockDevice) error {
	mkfsArgs := []string{
		"-I",
		"256",
		"-E",
		"lazy_itable_init=0,lazy_journal_init=0",
		device.Path(),
	}

	_, err := util.ExecuteCommand("mkfs.ext4", mkfsArgs...)
	return err
}

func resize2fs(device blockDevice) error {
	_, _ = util.ExecuteCommand("e2fsck", "-pf", device.Path())
	_, err := util.ExecuteCommand("resize2fs", device.Path())
	return err
}

func activateBackingDevice(file string, readOnly bool) (blockDevice, error) {
	fi, err := os.Stat(file)
	if err != nil {
		return nil, err
	}

	var device blockDevice

	if fi.Mode().IsRegular() {
		device = NewLoopDevice(file, readOnly)
	} else {
		// TODO: Support readOnly with physical devices somehow?
		device = newPhysicalDevice(file)
	}

	return device, device.activate()
}
