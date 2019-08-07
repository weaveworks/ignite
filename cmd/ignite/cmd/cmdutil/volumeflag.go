package cmdutil

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
)

var (
	formatErr = fmt.Errorf("volumes must be specified in the /host/path:/vm/path format")
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

		if len(paths) != 2 {
			return formatErr
		}

		volumeName := fmt.Sprintf("volume%d", i)

		// Create the Volume
		storage.Volumes = append(storage.Volumes, api.Volume{
			Name: volumeName,
			BlockDevice: &api.BlockDeviceVolume{
				Path: paths[0],
			},
		})

		// Create the VolumeMount
		storage.VolumeMounts = append(storage.VolumeMounts, api.VolumeMount{
			Name:      volumeName,
			MountPath: paths[1],
		})
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
