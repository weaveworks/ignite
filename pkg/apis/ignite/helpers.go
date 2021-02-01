package ignite

import (
	"path"

	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/util"
)

// SetImage populates relevant fields to an Image on the VM object
func (vm *VM) SetImage(image *Image) {
	vm.Spec.Image.OCI = image.Spec.OCI
	vm.Status.Image = image.Status.OCISource
}

// SetKernel populates relevant fields to a Kernel on the VM object
func (vm *VM) SetKernel(kernel *Kernel) {
	vm.Spec.Kernel.OCI = kernel.Spec.OCI
	vm.Status.Kernel = kernel.Status.OCISource
}

// NewPrefixer returns a util.Prefixer specific to the VM
func (vm *VM) NewPrefixer() *util.Prefixer {
	if vm.Status.IDPrefix == "" {
		// default for backwards compatibility /w previous VM's with no status set
		// note, we need to be careful that the VM.Status.IDPrefix is set as soon as possible
		// to avoid non-matching prefixedID's
		return util.NewPrefixer(constants.IGNITE_PREFIX)
	}
	return util.NewPrefixer(vm.Status.IDPrefix)
}

// PrefixedID returns the VM's prefixed ID for system references
func (vm *VM) PrefixedID() string {
	return vm.NewPrefixer().Prefix(vm.GetUID())
}

// SnapshotDev returns the path where the (legacy) DM snapshot exists
func (vm *VM) SnapshotDev() string {
	return path.Join("/dev/mapper", vm.PrefixedID())
}

// Running returns true if the VM is running, otherwise false
func (vm *VM) Running() bool {
	return vm.Status.Running
}

// OverlayFile returns the path to the overlay.dm file for the VM.
// TODO: This will be removed once we have the new snapshotter in place.
func (vm *VM) OverlayFile() string {
	return path.Join(vm.ObjectPath(), constants.OVERLAY_FILE)
}

// ObjectPath returns the directory where this VM's data is stored
func (vm *VM) ObjectPath() string {
	// TODO: Move this into storage
	return path.Join(constants.DATA_DIR, vm.GetKind().Lower(), vm.GetUID().String())
}

// ObjectPath returns the directory where this Image's data is stored
func (img *Image) ObjectPath() string {
	// TODO: Move this into storage
	return path.Join(constants.DATA_DIR, img.GetKind().Lower(), img.GetUID().String())
}

// ObjectPath returns the directory where this Kernel's data is stored
func (k *Kernel) ObjectPath() string {
	// TODO: Move this into storage
	return path.Join(constants.DATA_DIR, k.GetKind().Lower(), k.GetUID().String())
}
