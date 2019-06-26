package imgmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/dm"
	"github.com/weaveworks/ignite/pkg/format"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/source"
	"github.com/weaveworks/ignite/pkg/util"
)

var (
	dataDevSize = format.DataFrom(100 * 1073741824) // 100 GB
	blockSize   = format.DataFrom(65536)            // Pool allocation block size 128 (* 512 = 65536)
	extraSize   = format.DataFrom(100 * 1048576)    // Additional space to add to the image for the ext4 partition (100 MB)
)

func (md *ImageMetadata) NewDMPool() error {
	metadataFile := path.Join(md.ObjectPath(), constants.IMAGE_THINMETADATA)
	dataFile := path.Join(md.ObjectPath(), constants.IMAGE_THINDATA)

	if err := allocateDMFiles(metadataFile, dataFile); err != nil {
		return err
	}

	md.ImageOD().Pool = dm.NewPool(
		util.NewPrefixer().Prefix("pool", md.ID.String()),
		dataDevSize,
		blockSize,
		dm.NewLoopDevice(metadataFile, false),
		dm.NewLoopDevice(dataFile, false),
	)

	return nil
}

// Allocate the thin provisioning data and metadata files
func allocateDMFiles(metadataFile, dataFile string) error {
	thinFiles := map[string]format.DataSize{
		metadataFile: calcMetadataDevSize(dataDevSize),
		dataFile:     dataDevSize,
	}

	for p, size := range thinFiles {
		if !util.FileExists(p) {
			file, err := os.Create(p)
			if err != nil {
				return fmt.Errorf("failed to create thin provisioning file %q: %v", p, err)
			}

			// Allocate the image file
			if err := file.Truncate(int64(size.Bytes())); err != nil {
				return fmt.Errorf("failed to allocate space for thin provisioning file %q: %v", p, err)
			}

			if err := file.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (md *ImageMetadata) newImageVolume(src source.Source) (*util.MountPoint, error) {
	od := md.ImageOD()
	p := util.NewPrefixer()

	volume, err := od.Pool.CreateVolume(p.Prefix(src.ID()), src.Size().Add(extraSize))
	if err != nil {
		return nil, err
	}

	mountPoint, err := volume.Import(src)
	if err != nil {
		return nil, err
	}

	return mountPoint, nil
}

func (md *ImageMetadata) CreateOverlay(kernelSrc source.Source, requestedSize format.DataSize, id *metadata.ID) (*dm.Device, error) {
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

func calcMetadataDevSize(dataDevSize format.DataSize) format.DataSize {
	// The minimum size is 2 MB and the maximum size is 16 GB
	minSize := format.DataFrom(2 * 1048576)
	maxSize := format.DataFrom(16 * 1073741824)

	return format.DataFrom(48 * dataDevSize.Bytes() / blockSize.Bytes()).Min(maxSize).Max(minSize)
}
