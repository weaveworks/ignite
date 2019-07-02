package snapshotter

import (
	"fmt"
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"

	ignitemeta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

// Compile-time assert to verify interface compatibility
var _ Filter = &kernelFilter{}

// The KernelFilter filters kernels that belong to an image and have a specific size
type kernelFilter struct {
	image *Image
	size  ignitemeta.Size
}

func newKernelFilter(image *Image, size ignitemeta.Size) *kernelFilter {
	return &kernelFilter{
		image: image,
		size:  size,
	}
}

func (f *kernelFilter) SetType(t v1alpha1.PoolDeviceType) {}

func (f *kernelFilter) Filter(o *Object) (*Object, error) {
	mo, err := o.GetMetaObject()
	if err != nil {
		return nil, err
	}

	kernel, ok := mo.(*v1alpha1.Kernel)
	if !ok {
		return nil, fmt.Errorf("invalid object type for KernelFilter: %T", mo)
	}

	// Check the size
	if kernel.Spec.Source.Size != f.size {
		return nil, nil
	}

	// Check if child of image
	if !o.ChildOf(f.image) {
		return nil, nil
	}

	return o, nil
}

func (f *kernelFilter) ErrAmbiguous() ErrAmbiguous {
	return fmt.Errorf("kernelFilter: ambiguous")
}

func (f *kernelFilter) ErrNonexistent() ErrNonexistent {
	return fmt.Errorf("kernelFilter: nonexistent")
}
