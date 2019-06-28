package snapshotter

import (
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/dm"
	"github.com/weaveworks/ignite/pkg/layer"
	"github.com/weaveworks/ignite/pkg/source"
	"github.com/weaveworks/ignite/pkg/util"
	"image"
)

// This package represents image objects, which reside in /var/lib/firecracker/image/{id}/metadata.json
type Image struct {
	v1alpha1.Image
	layer
}

func NewImage(image v1alpha1.Image, device *dm.Device) *Image {
	return &Image{
		Image: image,
		layer: newLayer(device),
	}
}

func (i *Image) Load() {}

func (i *Image) Import() (*util.MountPoint, error) {
	volume, err := i.device.CreateVolume(i)
	if err != nil {
		return nil, err
	}

	// TODO: Different sources, for now just Docker
	src := source.NewDockerSource()
	if err := src.Parse(&image.Spec.Source); err != nil {
		return nil, err
	}

	mountPoint, err := volume.Import(src)
	if err != nil {
		return nil, err
	}

	return mountPoint, nil
}

var _ layer.Layer = &Image{}

func (i *Image) ID() *v1alpha1.DMID {
	return &i.Status.LayerID
}

// Get the metadata filename for the image
func (i *Image) MetadataPath() string {
	// TODO: This
	return ""
}

// TODO: This is the wrong size
func (i *Image) Size() v1alpha1.Size {
	return i.Spec.Source.Size
}
