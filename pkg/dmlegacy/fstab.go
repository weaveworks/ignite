package dmlegacy

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"text/tabwriter"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/util"
)

const (
	mountOptions = "rw,relatime"
)

var (
	blkidUUIDRegex = regexp.MustCompile("UUID=\"([^ ]*)\"")
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

func getUUID(devPath string) (string, error) {
	// running blkid requires root
	// we parse the output with regex because the `-o value -s UUID` format flags are not portable (ex: Alpine Linux)
	cmd := exec.Command("blkid", devPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command %q exited with %q: %w", cmd.Args, out, err)
	}

	uuidMatch := blkidUUIDRegex.FindStringSubmatch(string(out))
	if len(uuidMatch) > 1 {
		return uuidMatch[1], nil
	}

	return "", fmt.Errorf("command %q with output %q did not return a disk UUID", cmd.Args, out)
}
