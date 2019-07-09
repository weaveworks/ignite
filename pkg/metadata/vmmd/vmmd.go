package vmmd

import (
	"fmt"
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/filter"
	"github.com/weaveworks/ignite/pkg/metadata"
)

type VM struct {
	*api.VM
	// kernelUID and imageUID reference dependencies of the VM
	kernelUID meta.UID
	imageUID  meta.UID

	c *client.Client
}

var _ metadata.Metadata = &VM{}

// WrapVM wraps an API type in the runtime object
// It does not do any validation or checking like
// NewVM, hence it should only be used for "safe"
// data coming from storage.
func WrapVM(obj *api.VM) *VM {
	// Run the object through defaulting, just to be sure it has all the values
	scheme.Scheme.Default(obj)

	vm := &VM{
		VM: obj,
		c:  client.DefaultClient,
	}
	return vm
}

func NewVM(obj *api.VM, c *client.Client) (*VM, error) {
	// Initialize UID, name, defaulting, etc. that is common for all kinds
	if err := metadata.InitObject(obj, c); err != nil {
		return nil, err
	}

	// TODO: Validate the API object here

	// Construct the runtime object
	vm := &VM{
		VM: obj,
		c:  c,
	}

	// Populate dependent UIDs
	if err := vm.setImageUID(); err != nil {
		return nil, err
	}
	if err := vm.setKernelUID(); err != nil {
		return nil, err
	}
	return vm, nil
}

func (vm *VM) setImageUID() error {
	// TODO: Centralize validation
	if vm.Spec.Image.OCIClaim.Ref.IsUnset() {
		return fmt.Errorf("the image's OCIClaim ref field is mandatory")
	}

	image, err := vm.c.Images().Find(filter.NewNameFilter(vm.Spec.Image.OCIClaim.Ref.String()))
	if err != nil {
		return err
	}

	vm.imageUID = image.GetUID()
	return nil
}

func (vm *VM) setKernelUID() error {
	if vm.Spec.Kernel.OCIClaim.Ref.IsUnset() {
		return fmt.Errorf("the kernel's OCIClaim ref field is mandatory")
	}

	kernel, err := vm.c.Kernels().Find(filter.NewNameFilter(vm.Spec.Kernel.OCIClaim.Ref.String()))
	if err != nil {
		return err
	}

	vm.kernelUID = kernel.GetUID()
	return nil
}

func (vm *VM) GetImageUID() meta.UID {
	if len(vm.imageUID) == 0 {
		if err := vm.setImageUID(); err != nil {
			log.Debugf("Could not get image which this VM refers to: %v", err)
		}
	}
	return vm.imageUID
}

func (vm *VM) GetKernelUID() meta.UID {
	if len(vm.kernelUID) == 0 {
		if err := vm.setKernelUID(); err != nil {
			log.Debugf("Could not get kernel which this VM refers to: %v", err)
		}
	}
	return vm.kernelUID
}

func (vm *VM) ObjectPath() string {
	// TODO: Move this into storage
	return path.Join(constants.DATA_DIR, vm.GetKind().Lower(), vm.GetUID().String())
}

func (vm *VM) Save() error {
	return vm.c.VMs().Set(vm.VM)
}
