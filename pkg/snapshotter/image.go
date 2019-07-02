package snapshotter

import (
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	ignitemeta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/dm"
	"github.com/weaveworks/ignite/pkg/source"
	"github.com/weaveworks/ignite/pkg/util"
	"log"
)

// This package represents image objects, which reside in /var/lib/firecracker/image/{id}/metadata.json
type Image struct {
	*v1alpha1.Image
	device *dm.Device
}

func newImage(o *Object) (*Image, error) {
	mo, err := o.GetMetaObject()
	if err != nil {
		return nil, err
	}

	return &Image{
		Image:  mo.(*v1alpha1.Image),
		device: o.device,
	}, nil
}

// TODO: This should trigger the generation of a new image
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

func (i *Image) createResize(size ignitemeta.Size) (*Resize, error) {
	device, err := i.device.CreateSnapshot(size, "")
	if err != nil {
		return nil, err
	}

	// TODO: How to make this clean?
	// TODO: This is now not saved in the S
	return &Resize{
		dev:   device,
		image: i,
	}, nil
}

func (s *Snapshotter) CreateVM(image *Image, kernel *v1alpha1.Kernel, kernelSrc source.Source, vm *v1alpha1.VM) (*VM, error) {
	size := vm.Spec.Size.Max(image.device.Size) // The size needs to be at least the size of the image volume
	if size != vm.Spec.Size {
		// TODO: Issue a warning
		log.Printf("VM size %s < image size %s, using image size for VM", vm.Spec.Size.HR(), image.device.Size.HR())
	}

	var err error
	var resize *Resize

	if resize, err = s.GetResize(newResizeFilter(image, size)); err != nil {
		switch err.(type) {
		case ErrNonexistent:
			resize, err = image.device.crea(*v1alpha1.Kernel)
		default:
			return nil, err
		}
	}
}

func (s *Snapshotter) ImportKernel(image *Image, kernel *Kernel) error {
	volume := s.GetDevice(image.Status.LayerID)

	s.genResizeLayer(image)

	var err error
	p := util.NewPrefixer()
	pool := md.ImageOD().Pool

	volume, err := pool.Get(p.Prefix(md.ID.String()))
	if err != nil {
		return nil, err
	}

	// Make sure the overlay is always larger than the image
	// We need to do this here, as the size is used to name
	// the resize layers and everything on top
	size := requestedSize.Max(volume.Size())

	if size != requestedSize {
		// TODO: Warning error level
		log.Printf("Requested size %s < image size %s, using image size for overlay", requestedSize.HR(), size.HR())
	}

	resizeName := p.Prefix("resize", size.String())
	kernelName := p.Prefix("kernel")
	overlayName := p.Prefix(id.String())

	var kernel *dm.Device
	if kernel, err = pool.Get(kernelName); err != nil {
		// Requested kernel doesn't exist, so import it
		resize, err := volume.CreateSnapshot(resizeName, size)
		if err != nil {
			return nil, err
		}

		kernel, err = resize.CreateSnapshot(kernelName, resize.Size())
		if err != nil {
			return nil, err
		}

		mountPoint, err := kernel.Import(kernelSrc)
		if err != nil {
			return nil, err
		}

		err = mountPoint.Umount()
		if err != nil {
			return nil, err
		}
	}

	overlay, err := kernel.CreateSnapshot(overlayName, kernel.Size())
	if err != nil {
		return nil, err
	}

	return overlay, nil
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
