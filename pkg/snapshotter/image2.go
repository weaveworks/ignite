package snapshotter

import (
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/dm"
	"github.com/weaveworks/ignite/pkg/layer/image"
	"github.com/weaveworks/ignite/pkg/layer/kernel"
	"github.com/weaveworks/ignite/pkg/layer/resize"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/source"
	"github.com/weaveworks/ignite/pkg/util"
	"log"
)

func (s *Snapshotter) genResizeLayer(image *image.Image, size v1alpha1.Size) (*dm.Device, error) {
	for i, device := range s.Status.Devices {
		// TODO: Improve this test, make sure it matches only resize layers
		if device != nil && device.Parent == image.Status.LayerID && device.Size == size && len(device.MetadataPath) == 0 {
			return s.GetDevice(v1alpha1.NewDMID(i)), nil
		}
	}

	return s.GetDevice(*image.ID()).CreateSnapshot(resize.NewResize(size))
}

func (s *Snapshotter) ImportKernel(image *image.Image, kernel *kernel.Kernel) error {
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
