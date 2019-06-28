package snapshotter

import (
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/dm"
)

type Resize struct {
	layer
}

func NewResize(device *dm.Device) *Resize {
	return &Resize{
		newLayer(device),
	}
}

func (r *Resize) ID() *v1alpha1.DMID {
	return &r.layerID
}

// Resize layers have no metadata
func (r *Resize) MetadataPath() string {
	return ""
}

func (r *Resize) Size() v1alpha1.Size {
	return r.size
}
