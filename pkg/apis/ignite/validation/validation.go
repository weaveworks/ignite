package validation

import (
	"fmt"
	"path"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateVM validates a VM object and collects all encountered errors
func ValidateVM(obj *api.VM) (allErrs field.ErrorList) {
	allErrs = append(allErrs, ValidateNetworkMode(obj.Spec.Network.Mode, field.NewPath(".spec.network.mode"))...)
	allErrs = append(allErrs, ValidateOCIImageClaim(&obj.Spec.Image.OCIClaim, field.NewPath(".spec.image.ociClaim"))...)
	allErrs = append(allErrs, ValidateOCIImageClaim(&obj.Spec.Kernel.OCIClaim, field.NewPath(".spec.kernel.ociClaim"))...)
	allErrs = append(allErrs, ValidateFileMappings(&obj.Spec.CopyFiles, field.NewPath(".spec.copyFiles"))...)
	allErrs = append(allErrs, ValidateVMState(obj.Status.State, field.NewPath(".status.state"))...)
	// TODO: Add vCPU, memory, disk max and min sizes
	return
}

// ValidateOCIImageClaim validates an OCI image claim
func ValidateOCIImageClaim(c *api.OCIImageClaim, fldPath *field.Path) (allErrs field.ErrorList) {
	allErrs = append(allErrs, ValidateImageSourceType(c.Type, fldPath.Child("type"))...)
	return
}

// ValidateFileMappings validates if the filemappings is valid
func ValidateFileMappings(mappings *[]api.FileMapping, fldPath *field.Path) (allErrs field.ErrorList) {
	for i, mapping := range *mappings {
		mappingPath := fldPath.Child(fmt.Sprintf("[%d]", i))
		allErrs = append(allErrs, ValidateAbsolutePath(mapping.HostPath, mappingPath.Child("hostPath"))...)
		allErrs = append(allErrs, ValidateAbsolutePath(mapping.VMPath, mappingPath.Child("vmPath"))...)
	}
	return
}

// ValidateAbsolutePath validates if a path is absolute
func ValidateAbsolutePath(pathStr string, fldPath *field.Path) (allErrs field.ErrorList) {
	if !path.IsAbs(pathStr) {
		allErrs = append(allErrs, field.Invalid(fldPath, pathStr, fmt.Sprintf("path must be absolute %q", pathStr)))
	}
	return
}

// ValidateNetworkMode validates if a network mode is valid
func ValidateNetworkMode(mode api.NetworkMode, fldPath *field.Path) (allErrs field.ErrorList) {
	found := false
	modes := api.GetNetworkModes()
	for _, nm := range modes {
		if nm == mode {
			found = true
		}
	}
	if !found {
		allErrs = append(allErrs, field.Invalid(fldPath, mode, fmt.Sprintf("network mode must be one of %v", modes)))
	}
	return
}

// ValidateImageSourceType validates if an image source type is valid
func ValidateImageSourceType(t api.ImageSourceType, fldPath *field.Path) (allErrs field.ErrorList) {
	found := false
	types := api.GetImageSourceTypes()
	for _, tt := range types {
		if tt == t {
			found = true
		}
	}
	if !found {
		allErrs = append(allErrs, field.Invalid(fldPath, t, fmt.Sprintf("image source type must be one of %v", types)))
	}
	return
}

// ValidateVMState validates if an VM state is valid
func ValidateVMState(s api.VMState, fldPath *field.Path) (allErrs field.ErrorList) {
	found := false
	states := api.GetVMStates()
	for _, state := range states {
		if state == s {
			found = true
		}
	}
	if !found {
		allErrs = append(allErrs, field.Invalid(fldPath, s, fmt.Sprintf("VM state must be one of %v", states)))
	}
	return
}
