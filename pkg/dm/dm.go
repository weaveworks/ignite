package dm

import (
	"fmt"
	"log"
	"os"

	"github.com/weaveworks/ignite/pkg/util"

	ignitemeta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
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

func allocateBackingFile(p string, size ignitemeta.Size) error {
	if !util.FileExists(p) {
		file, err := os.Create(p)
		if err != nil {
			return fmt.Errorf("failed to create thin provisioning file %q: %v", p, err)
		}

		// Allocate the image file
		if err := file.Truncate(int64(size.Bytes())); err != nil {
			return fmt.Errorf("failed to allocate space for thin provisioning file %q: %v", p, err)
		}

		if err := file.Close(); err != nil {
			return err
		}
	}

	return nil
}

func activateBackingDevice(p string, readOnly bool) (blockDevice, error) {
	fi, err := os.Stat(p)
	if err != nil {
		return nil, err
	}

	var device blockDevice

	if fi.Mode().IsRegular() {
		device = NewLoopDevice(p, readOnly)
	} else {
		// TODO: Support readOnly with physical devices somehow?
		device = newPhysicalDevice(p)
	}

	return device, device.activate()
}
