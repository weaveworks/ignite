package dmlegacy

import (
	"fmt"
	"os"
	"path"

	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/util"
)

// ActivateSnapshot sets up the snapshot with devicemapper so that it is active and can be used
func ActivateSnapshot(vm *vmmd.VM) error {
	device := constants.IGNITE_PREFIX + vm.GetUID().String()
	devicePath := vm.SnapshotDev()

	// Return if the snapshot is already setup
	if util.FileExists(devicePath) {
		return nil
	}

	// Setup loop device for the image
	imageLoop, err := newLoopDev(path.Join(constants.IMAGE_DIR, vm.GetImageUID().String(), constants.IMAGE_FS), true)
	if err != nil {
		return err
	}

	// Make sure the all directories above the snapshot directory exists
	if err := os.MkdirAll(path.Dir(vm.OverlayFile()), 0755); err != nil {
		return err
	}

	// Setup loop device for the VM overlay
	overlayLoop, err := newLoopDev(vm.OverlayFile(), false)
	if err != nil {
		return err
	}

	imageLoopSize, err := imageLoop.Size512K()
	if err != nil {
		return err
	}

	overlayLoopSize, err := overlayLoop.Size512K()
	if err != nil {
		return err
	}

	// If the overlay is larger than the base image, we need to set up an additional dm device
	// which will contain the image and additional zero space (which reads zeros and discards writes).
	// This is fine, because all writes will target the overlay snapshot and not the read-only image.
	// The newly generated larger device will then be used for creating the snapshot (which is always
	// as large as the device backing it).

	basePath := imageLoop.Path()
	if overlayLoopSize > imageLoopSize {
		// "0 8388608 linear /dev/loop0 0"
		// "8388608 12582912 zero"
		dmBaseTable := []byte(fmt.Sprintf("0 %d linear %s 0\n%d %d zero", imageLoopSize, imageLoop.Path(), imageLoopSize, overlayLoopSize))

		baseDevice := fmt.Sprintf("%s-base", device)
		if err := runDMSetup(baseDevice, dmBaseTable); err != nil {
			return err
		}

		basePath = fmt.Sprintf("/dev/mapper/%s", baseDevice)
	}

	// "0 8388608 snapshot /dev/{loop0,mapper/ignite-<uid>-base} /dev/loop1 P 8"
	dmTable := []byte(fmt.Sprintf("0 %d snapshot %s %s P 8", overlayLoopSize, basePath, overlayLoop.Path()))

	if err := runDMSetup(device, dmTable); err != nil {
		return err
	}

	// Repair the filesystem in case it has errors
	// e2fsck throws an error if the filesystem gets repaired, so just ignore it
	_, _ = util.ExecuteCommand("e2fsck", "-p", "-f", devicePath)

	// If the overlay is larger than the image, call resize2fs to make the filesystem fill the overlay
	if overlayLoopSize > imageLoopSize {
		if _, err := util.ExecuteCommand("resize2fs", devicePath); err != nil {
			return err
		}
	}

	// By detaching the loop devices after setting up the snapshot
	// they get automatically removed when the snapshot is removed.
	if err := imageLoop.Detach(); err != nil {
		return err
	}

	return overlayLoop.Detach()
}

// DeactivateSnapshot deactivates the snapshot by removing it with dmsetup
func DeactivateSnapshot(vm *vmmd.VM) error {
	dmArgs := []string{
		"remove",
		vm.SnapshotDev(),
	}

	// If the base device is visible in "dmsetup", we should remove it
	// The device itself is not forwarded to docker, so we can't query its path
	// TODO: Improve this detection
	baseDev := fmt.Sprintf("%s-base", constants.IGNITE_PREFIX+vm.GetUID())
	if _, err := util.ExecuteCommand("dmsetup", "info", baseDev); err == nil {
		dmArgs = append(dmArgs, baseDev)
	}

	_, err := util.ExecuteCommand("dmsetup", dmArgs...)
	return err
}
