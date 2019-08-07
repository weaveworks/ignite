package dmlegacy

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/tabwriter"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/util"
)

const (
	mountOptions = "rw,relatime"
	uuidPath     = "/dev/disk/by-uuid"
)

type fstabEntry struct {
	uuid       string
	mountPoint string
}

var _ fmt.Stringer = &fstabEntry{}

func (f *fstabEntry) isValid() bool {
	// An entry is valid if both the UUID and mount point are set
	return len(f.uuid) > 0 && len(f.mountPoint) > 0
}

func (f *fstabEntry) String() string {
	return strings.Join([]string{
		fmt.Sprintf("UUID=%s", f.uuid), // Mount by UUID
		f.mountPoint,                   // The mount point for the volume
		"auto",                         // Discover the filesystem automatically
		mountOptions,                   // Use the mount options defined above
		"0",                            // Don't dump the filesystem
		"2",                            // fsck may check this filesystem on reboot
	}, "\t")
}

func populateFstab(vm *api.VM, mountPoint string) error {
	fstab, err := os.Create(path.Join(mountPoint, "/etc/fstab"))
	if err != nil {
		return err
	}
	defer fstab.Close()

	writer := new(tabwriter.Writer)
	writer.Init(fstab, 0, 8, 1, '\t', 0)
	entries := make(map[string]*fstabEntry, util.MaxInt(len(vm.Spec.Storage.Volumes), len(vm.Spec.Storage.VolumeMounts)))

	// Discover all volumes
	for _, volume := range vm.Spec.Storage.Volumes {
		if volume.BlockDevice == nil {
			continue // Skip all non block device volumes for now
		}

		// Retrieve the UUID for the block device
		uuid, err := getUUID(volume.BlockDevice.Path)
		if err != nil {
			return err
		}

		// Create an entry for the volume
		entries[volume.Name] = &fstabEntry{uuid: uuid}
	}

	// Discover all volume mounts
	for _, volumeMount := range vm.Spec.Storage.VolumeMounts {
		// Lookup the entry based on the volume mount's name
		if entry, ok := entries[volumeMount.Name]; ok {
			// Add the mount path to the entry if it exists
			entry.mountPoint = volumeMount.MountPath
		}
	}

	for _, entry := range entries {
		if entry.isValid() {
			// Write the entry to /etc/fstab
			if _, err := fmt.Fprint(writer, entry, "\n"); err != nil {
				return err
			}
		} else {
			// This should have been caught in validation
			return fmt.Errorf("invalid fstab entry: %q -> %q", entry.uuid, entry.mountPoint)
		}
	}

	return writer.Flush()
}

func getUUID(devPath string) (uuid string, err error) {
	var files []os.FileInfo
	if files, err = ioutil.ReadDir(uuidPath); err != nil {
		return
	}

	for _, fi := range files {
		if fi.Mode()&os.ModeSymlink == 0 {
			continue // Skip all non-symlinks
		}

		// Resolve the symbolic link
		var link string
		if link, err = os.Readlink(path.Join(uuidPath, fi.Name())); err != nil {
			return
		}

		// Make the link target absolute and compare with the given path
		if path.Join(uuidPath, link) == devPath {
			uuid = fi.Name() // The UUID is the filename of the symbolic link
			break
		}
	}

	if len(uuid) == 0 {
		err = fmt.Errorf("no UUID found for device %q, verify it has a filesystem", devPath)
	}

	return
}
