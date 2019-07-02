package snapshotter

import (
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/dm"
)

type VM struct {
	*v1alpha1.VM
	device *dm.Device
	kernel *Kernel
}

func newVM(o *Object) (*VM, error) {
	mo, err := o.GetMetaObject()
	if err != nil {
		return nil, err
	}

	kernel, err := newKernel(o.parent)
	if err != nil {
		return nil, err
	}

	return &VM{
		VM:     mo.(*v1alpha1.VM),
		device: o.device,
		kernel: kernel,
	}, nil
}

func (s *Snapshotter) createVM(vm *v1alpha1.VM) (*VM, error) {
	kernelObj := &Object{}

	var err error
	var kernel *Kernel

	if kernel, err = s.GetKernel(newKernelFilter(kernelObj)); err != nil {
		switch err.(type) {
		case ErrNonexistent:
			kernel, err = s.createKernel(*v1alpha1.Kernel)
		default:
			return nil, err
		}
	}



	o := &Object{
		device: s.,
		object: vm,
		parent: nil,
	}
}