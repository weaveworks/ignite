package cmdutil

import (
	"fmt"
	"path"
	"strings"

	"github.com/spf13/pflag"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/util"
)

var (
	formatErr   = fmt.Errorf("volumes must be specified in the /host/path:/vm/path format")
	absPathsErr = fmt.Errorf("paths given to the volume flag must be absolute")
)

// VolumeFlag is the pflag.Value custom flag for `ignite create --volume`
type VolumeFlag struct {
	value *api.VMStorageSpec
	s     string
}

var _ pflag.Value = &VolumeFlag{}

func (vf *VolumeFlag) Set(x string) error {
	entries := strings.Split(x, ",") // We take in a comma-separated list
	var storage api.VMStorageSpec

	for i, entry := range entries {
		paths := strings.Split(entry, ":")

		var hostPath, vmPath string
		switch len(paths) {
		case 0:
			// No paths given
			continue
		case 2: // We were given
			// Set the VM path and verify it's absolute
			if vmPath = paths[1]; !path.IsAbs(vmPath) {
				return absPathsErr
			}

			fallthrough
		case 1:
			// Generate a dummy volume name
			volumeName := fmt.Sprintf("volume%d", i)

			// Set the host path and verify it's absolute
			if hostPath = paths[0]; !path.IsAbs(hostPath) {
				return absPathsErr
			}

			// Verify that the host path points to a device file
			if err := util.DeviceFile(hostPath); err != nil {
				return err
			}

			// If the VM path is not set, use the host path
			if len(vmPath) == 0 {
				vmPath = hostPath
			}

			// Create the Volume
			// TODO: Check for name collisions
			storage.Volumes = append(storage.Volumes, api.Volume{
				Name: volumeName,
				BlockDevice: &api.BlockDeviceVolume{
					Path: hostPath,
				},
			})

			// Create the VolumeMount
			// TODO: Check for name collisions
			storage.VolumeMounts = append(storage.VolumeMounts, api.VolumeMount{
				Name:      volumeName,
				MountPath: vmPath,
			})
		default:
			// Incorrect format
			return formatErr
		}
	}

	*vf.value = storage
	vf.s = x // String should return the input after a successful Set
	return nil
}

func (vf *VolumeFlag) Type() string {
	return "volume"
}

func (vf *VolumeFlag) String() string {
	return vf.s
}

func VolumeVar(fs *pflag.FlagSet, ptr *api.VMStorageSpec, name, usage string) {
	VolumeVarP(fs, ptr, name, "", usage)
}

func VolumeVarP(fs *pflag.FlagSet, ptr *api.VMStorageSpec, name, shorthand, usage string) {
	fs.VarP(&VolumeFlag{value: ptr}, name, shorthand, usage)
}
