package validation

import (
	"fmt"
	"path"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/util"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateVM validates a VM object and collects all encountered errors
func ValidateVM(obj *api.VM) (allErrs field.ErrorList) {
	allErrs = append(allErrs, ValidateVMName(obj.GetName(), field.NewPath("metadata.name"))...)
	allErrs = append(allErrs, RequireOCIImageRef(&obj.Spec.Image.OCI, field.NewPath(".spec.image.oci"))...)
	allErrs = append(allErrs, RequireOCIImageRef(&obj.Spec.Kernel.OCI, field.NewPath(".spec.kernel.oci"))...)
	allErrs = append(allErrs, ValidateFileMappings(&obj.Spec.CopyFiles, field.NewPath(".spec.copyFiles"))...)
	allErrs = append(allErrs, ValidateVMStorage(&obj.Spec.Storage, field.NewPath(".spec.storage"))...)
	// TODO: Add vCPU, memory, disk max and min sizes
	// TODO: Add port mapping validation
	return
}

// RequireOCIImageRef validates that the OCIImageRef is set
func RequireOCIImageRef(ref *meta.OCIImageRef, fldPath *field.Path) (allErrs field.ErrorList) {
	if ref.IsUnset() {
		allErrs = append(allErrs, field.Required(fldPath, "the OCI reference is mandatory"))
	}

	return
}

// ValidateVMName validates the VM name.
func ValidateVMName(name string, fldPath *field.Path) (allErrs field.ErrorList) {
	errs := validation.IsDNS1123Subdomain(name)
	for _, e := range errs {
		allErrs = append(allErrs, field.Invalid(fldPath, name, e))
	}

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

// ValidateNonemptyName validated that the given name is nonempty
func ValidateNonemptyName(name string, fldPath *field.Path) (allErrs field.ErrorList) {
	if util.IsEmptyString(name) {
		allErrs = append(allErrs, field.Invalid(fldPath, name, "name must be non-empty"))
	}

	return
}

// ValidateAbsolutePath validates if a path is absolute
func ValidateAbsolutePath(pathStr string, fldPath *field.Path) (allErrs field.ErrorList) {
	if !path.IsAbs(pathStr) {
		allErrs = append(allErrs, field.Invalid(fldPath, pathStr, "path must be absolute"))
	}

	return
}
