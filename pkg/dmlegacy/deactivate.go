package dmlegacy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nightlyone/lockfile"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/util"
)

// dmsetupNotFound is the error message when dmsetup can't find a device.
const dmsetupNotFound = "No such device or address"

// DeactivateSnapshot deactivates the snapshot by removing it with dmsetup
func DeactivateSnapshot(vm *api.VM) error {
	// Global lock path.
	glpath := filepath.Join(os.TempDir(), snapshotLockFileName)

	// Create a lockfile and obtain a lock.
	lock, err := lockfile.New(glpath)
	if err != nil {
		err = fmt.Errorf("failed to create lockfile: %w", err)
		return err
	}
	if err = obtainLock(lock); err != nil {
		return err
	}
	// Release the lock at the end.
	defer util.DeferErr(&err, lock.Unlock)

	dmArgs := []string{
		"remove",
		"--verifyudev", // if udevd is not running, dmsetup will manage the device node in /dev/mapper
		vm.PrefixedID(),
	}

	// If the base device is visible in "dmsetup", we should remove it
	// The device itself is not forwarded to docker, so we can't query its path
	// TODO: Improve this detection
	baseDev := vm.NewPrefixer().Prefix(vm.GetUID(), "base")
	if _, err := util.ExecuteCommand("dmsetup", "info", baseDev); err == nil {
		dmArgs = append(dmArgs, baseDev)
	}

	if _, err := util.ExecuteCommand("dmsetup", dmArgs...); err != nil {
		// If the device is not found, it's been deleted already, return nil.
		if strings.Contains(err.Error(), dmsetupNotFound) {
			return nil
		}
		return err
	}

	return nil
}
