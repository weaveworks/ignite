package snapshotter

import (
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/dm"
)

// This package represents kernel objects, which reside in /var/lib/firecracker/kernel/{id}/metadata.json
type Kernel struct {
	v1alpha1.Kernel
	layer
}

func NewKernel(kernel v1alpha1.Kernel, device *dm.Device) *Kernel {
	return &Kernel{
		Kernel: kernel,
		layer:  newLayer(device),
	}
}

// Get the metadata filename for the image
func (i *Image) MetadataPath() string {
	// TODO: This
	return ""
}
