package validation

import (
	"fmt"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/util"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateBlockDeviceVolume validates if the BlockDeviceVolume is valid
func ValidateBlockDeviceVolume(b *api.BlockDeviceVolume, fldPath *field.Path, paths map[string]struct{}) (allErrs field.ErrorList) {
	pathFldPath := fldPath.Child("path")
	allErrs = append(allErrs, ValidateAbsolutePath(b.Path, pathFldPath)...)

	// Validate that the block device path points to a device file
	if err := util.IsDeviceFile(b.Path); err != nil {
		allErrs = append(allErrs, field.Invalid(pathFldPath, b.Path, err.Error()))
	}

	// Validate path uniqueness
	if _, ok := paths[b.Path]; ok {
		allErrs = append(allErrs, field.Invalid(pathFldPath, b.Path, "blockDevice path must be unique"))
	} else {
		paths[b.Path] = struct{}{}
	}

	return
}

// ValidateVMStorage validates if the VMStorageSpec is valid
func ValidateVMStorage(s *api.VMStorageSpec, fldPath *field.Path) (allErrs field.ErrorList) {
	// names keeps track of volume names and if they have a respective volumeMount
	names := make(map[string]bool, util.MaxInt(len(s.VolumeMounts), len(s.Volumes)))
	// blockDevPaths keeps track of registered block device paths
	blockDevPaths := make(map[string]struct{}, len(s.Volumes))
	// mountPaths keeps track of registered volumeMount paths
	mountPaths := make(map[string]struct{}, len(s.VolumeMounts))

	// volume validation
	for i, volume := range s.Volumes {
		volumeFldPath := fldPath.Child(fmt.Sprintf("[%d]", i))
		allErrs = append(allErrs, ValidateNonemptyName(volume.Name, volumeFldPath.Child("name"))...)

		// For now require and validate the BlockDevice entry
		blockDevFldPath := volumeFldPath.Child("blockDevice")
		if volume.BlockDevice == nil {
			allErrs = append(allErrs, field.Invalid(blockDevFldPath, nil, "blockDevice must be non-nil"))
		} else {
			allErrs = append(allErrs, ValidateBlockDeviceVolume(volume.BlockDevice, blockDevFldPath, blockDevPaths)...)
		}

		// Validate volume name uniqueness
		if _, ok := names[volume.Name]; ok {
			allErrs = append(allErrs, field.Invalid(volumeFldPath.Child("name"), volume.Name, "volume name must be unique"))
		} else {
			names[volume.Name] = false
		}
	}

	// volumeMount validation
	for i, mount := range s.VolumeMounts {
		mountFldPath := fldPath.Child(fmt.Sprintf("[%d]", i))
		mountNameFldPath := mountFldPath.Child("name")
		mountPathFldPath := mountFldPath.Child("mountPath")

		allErrs = append(allErrs, ValidateNonemptyName(mount.Name, mountNameFldPath)...)
		allErrs = append(allErrs, ValidateAbsolutePath(mount.MountPath, mountPathFldPath)...)

		// Validate volumeMount name uniqueness and correlation to volumes
		if matched, ok := names[mount.Name]; ok {
			if matched {
				allErrs = append(allErrs, field.Invalid(mountNameFldPath, mount.Name, "volumeMount name must be unique"))
			} else {
				names[mount.Name] = true
			}
		} else {
			allErrs = append(allErrs, field.Invalid(mountNameFldPath, mount.Name, "volumeMount name must match a volume name"))
		}

		// Validate volumeMount path uniqueness
		if _, ok := mountPaths[mount.MountPath]; ok {
			allErrs = append(allErrs, field.Invalid(mountPathFldPath, mount.MountPath, "volumeMount path must be unique"))
		} else {
			mountPaths[mount.MountPath] = struct{}{}
		}
	}

	return
}
