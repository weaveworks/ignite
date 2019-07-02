package filter

import (
	"fmt"
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/snapshotter"

	ignitemeta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

// Compile-time assert to verify interface compatibility
var _ snapshotter.Filter = &KernelFilter{}

// The KernelFilter filters kernels that belong to an image and have a specific size
type KernelFilter struct {
	*IDNameFilter
	image *snapshotter.Image
	size  ignitemeta.Size
}

func NewKernelFilter(p string, image *snapshotter.Image, size ignitemeta.Size) *KernelFilter {
	return &KernelFilter{
		IDNameFilter: NewIDNameFilter(p),
		image:        image,
		size:         size,
	}
}

func (f *KernelFilter) Filter(object *snapshotter.Object) (*snapshotter.Object, error) {
	mo, err := object.GetMetaObject()
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
	if !object.ChildOf(f.image) {
		return nil, nil
	}

	return f.IDNameFilter.Filter(object)
}

func (f *KernelFilter) ErrAmbiguous() snapshotter.ErrAmbiguous {
	return fmt.Errorf("ambiguous %s query: %q matched the following IDs/names: %s", f.filterType, f.prefix, formatMatches(f.matches))
}

func (f *KernelFilter) ErrNonexistent() snapshotter.ErrNonexistent {
	return fmt.Errorf("can't find %s: no ID/name matches for %q", f.filterType, f.prefix)
}
