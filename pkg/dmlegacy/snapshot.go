package dmlegacy

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/nightlyone/lockfile"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/operations/lookup"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/util"
)

const snapshotLockFileName = "ignite-snapshot.lock"

// ActivateSnapshot sets up the snapshot with devicemapper so that it is active and can be used.
// It returns the path of the bootable snapshot device.
func ActivateSnapshot(vm *api.VM) (devicePath string, err error) {
	device := vm.PrefixedID()
	devicePath = vm.SnapshotDev()

	// Return if the snapshot is already setup
	if util.FileExists(devicePath) {
		return
	}

	// Get the image UID from the VM
	imageUID, err := lookup.ImageUIDForVM(vm, providers.Client)
	if err != nil {
		return
	}

	// NOTE: Multiple ignite processes trying to create loop devices at the
	// same time results in race condition. When multiple processes request for
	// a free loop device at the same time, they may get the same device ID and
	// try to create the same device multiple times.
	// Serialize this operation by creating a global lock file when creating a
	// loop device and release the lock after setting up device mapper using the
	// loop device.
	// Also, concurrent interactions with device mapper results in the error:
	//   Device or resource busy
	// Serializing the interaction with device mapper helps avoid issues due to
	// race conditions.

	// Global lock path.
	glpath := filepath.Join(os.TempDir(), snapshotLockFileName)

	// Create a lockfile and obtain a lock.
	lock, err := lockfile.New(glpath)
	if err != nil {
		err = fmt.Errorf("failed to create lockfile: %w", err)
		return
	}
	if err = obtainLock(lock); err != nil {
		return
	}
	// Release the lock at the end.
	defer util.DeferErr(&err, lock.Unlock)

	// Setup loop device for the image
	imageLoop, err := newLoopDev(path.Join(constants.IMAGE_DIR, imageUID.String(), constants.IMAGE_FS), true)
	if err != nil {
		return
	}

	// Make sure the all directories above the snapshot directory exists
	if err = os.MkdirAll(path.Dir(vm.OverlayFile()), 0755); err != nil {
		return
	}

	// Setup loop device for the VM overlay
	overlayLoop, err := newLoopDev(vm.OverlayFile(), false)
	if err != nil {
		return
	}

	imageLoopSize, err := imageLoop.Size512K()
	if err != nil {
		return
	}

	overlayLoopSize, err := overlayLoop.Size512K()
	if err != nil {
		return
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
		if err = runDMSetup(baseDevice, dmBaseTable); err != nil {
			return
		}

		basePath = fmt.Sprintf("/dev/mapper/%s", baseDevice)
	}

	// "0 8388608 snapshot /dev/{loop0,mapper/ignite-<uid>-base} /dev/loop1 P 8"
	dmTable := []byte(fmt.Sprintf("0 %d snapshot %s %s P 8", overlayLoopSize, basePath, overlayLoop.Path()))

	// setup the main boot device
	if err = runDMSetup(device, dmTable); err != nil {
		return
	}

	// Repair the filesystem in case it has errors
	// e2fsck throws an error if the filesystem gets repaired, so just ignore it
	_, _ = util.ExecuteCommand("e2fsck", "-p", "-f", devicePath)

	// If the overlay is larger than the image, call resize2fs to make the filesystem fill the overlay
	if overlayLoopSize > imageLoopSize {
		if _, err = util.ExecuteCommand("resize2fs", devicePath); err != nil {
			return
		}
	}

	// By detaching the loop devices after setting up the snapshot
	// they get automatically removed when the snapshot is removed.
	if err = imageLoop.Detach(); err != nil {
		return
	}

	err = overlayLoop.Detach()

	return
}

// obtainLock tries to obtain a lock and retries if the lock is owned by
// another process, until a lock is obtained.
func obtainLock(lock lockfile.Lockfile) error {
	// Check if the lock has any owner.
	process, err := lock.GetOwner()
	if err == nil {
		// A lock already exists. Check if the lock owner is the current process
		// itself.
		if process.Pid == os.Getpid() {
			return fmt.Errorf("lockfile %q already locked by this process", lock)
		}

		// A lock already exists, but it's owned by some other process. Continue
		// to obtain lock, in case the lock owner no longer exists.
	}

	// Obtain a lock. Retry if the lock can't be obtained.
	err = lock.TryLock()
	for err != nil {
		// Check if it's a lock temporary error that can be mitigated with a
		// retry. Fail if any other error.
		if _, ok := err.(interface{ Temporary() bool }); !ok {
			return fmt.Errorf("unable to lock %q: %v", lock, err)
		}
		err = lock.TryLock()
	}

	return nil
}
