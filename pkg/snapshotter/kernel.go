package snapshotter

import (
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/dm"
)

// This package represents kernel objects, which reside in /var/lib/firecracker/kernel/{id}/metadata.json
type Kernel struct {
	*Object
	*v1alpha1.Kernel
	resize *Resize
}

func newKernel(o *Object) (*Kernel, error) {
	mo, err := o.GetMetaObject()
	if err != nil {
		return nil, err
	}

	resize, err := newResize(o.parent)
	if err != nil {
		return nil, err
	}

	return &Kernel{
		Kernel: mo.(*v1alpha1.Kernel),
		device: o.device,
		resize: resize,
	}, nil
}

func (k *Kernel) resize() *Resize {
	return k.parent.(*Resize)
}

func (s *Snapshotter) createKernel() (*Kernel, error) {
	o := &Object{
		device: s.,
		object: nil,
		parent: nil,
	}
}

func (k *Kernel) ChildOf(image *Image) bool {
	return k.resize.ChildOf(image)
}
