package snapshotter

import (
	"fmt"
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	ignitemeta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/dm"
)

type Resize struct {
	dev   *dm.Device
	image *Image
}

var _ Object = &Resize{}

func newResize(o *Object) (*Resize, error) {
	image, err := newImage(o.parent)
	if err != nil {
		return nil, err
	}

	return &Resize{
		dev:   o.device,
		image: image,
	}, nil
}

// TODO: This should generate a new resize from the ground up
func WholeNewResize() {

}

func (r *Resize) device() *dm.Device {
	return r.dev
}

func (r *Resize) ChildOf(image *Image) bool {
	// The dm.Devices are canonically saved in Snapshotter
	return image.device == r.image.device
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

// Compile-time assert to verify interface compatibility
var _ Filter = &resizeFilter{}

// The resizeFilter filters specific resize objects
type resizeFilter struct {
	image *Image
	size  ignitemeta.Size
}

func newResizeFilter(image *Image, size ignitemeta.Size) *resizeFilter {
	return &resizeFilter{
		image: image,
		size:  size,
	}
}

func (f *resizeFilter) SetType(t v1alpha1.PoolDeviceType) {}

func (f *resizeFilter) Filter(o *Object) (*Object, error) {
	// Check the size
	if o.device.Size != f.size {
		return nil, nil
	}

	// Check if child of image
	if !o.ChildOf(f.image) {
		return nil, nil
	}

	return o, nil
}

func (f *resizeFilter) ErrAmbiguous() ErrAmbiguous {
	return fmt.Errorf("resizeFilter: ambiguous")
}

func (f *resizeFilter) ErrNonexistent() ErrNonexistent {
	return fmt.Errorf("resizeFilter: nonexistent")
}
