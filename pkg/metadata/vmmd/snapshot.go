package vmmd

import (
	"fmt"
	"github.com/freddierice/go-losetup"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/util"
	"path"
	"strings"
)

func (md *VMMetadata) SetupSnapshot() (string, error) {
	device := constants.IGNITE_PREFIX + md.ID
	devicePath := path.Join("/dev/mapper", device)

	// Return if the snapshot is already setup
	if util.FileExists(devicePath) {
		return devicePath, nil
	}

	// Setup loop device for the image
	imageLoop, err := newLoopDev(path.Join(constants.IMAGE_DIR, md.VMOD().ImageID, constants.IMAGE_FS), true)
	if err != nil {
		return "", err
	}

	// Setup loop device for the VM overlay
	overlayLoop, err := newLoopDev(path.Join(md.ObjectPath(), constants.OVERLAY_FILE), false)
	if err != nil {
		return "", err
	}

	imageLoopSize, err := util.ExecuteCommand("blockdev", "--getsz", imageLoop.Path())
	if err != nil {
		return "", err
	}

	// dmsetup create newdev --table "0 8388608 snapshot /dev/loop0 /dev/loop1 P 8"
	dmTable := []string{
		"0",
		imageLoopSize,
		"snapshot",
		imageLoop.Path(),
		overlayLoop.Path(),
		"P",
		"8",
	}

	dmArgs := []string{
		"create",
		device,
		"--table",
		strings.Join(dmTable, " "),
	}

	if _, err := util.ExecuteCommand("dmsetup", dmArgs...); err != nil {
		return "", err
	}

	return devicePath, nil
}

func newLoopDev(file string, readOnly bool) (*losetup.Device, error) {
	dev, err := losetup.Attach(file, 0, readOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to setup loop device for %q: %v", file, err)
	}

	return &dev, nil
}
