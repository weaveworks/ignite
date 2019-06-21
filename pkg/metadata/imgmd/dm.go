package imgmd

import (
	"fmt"
	"os"
	"path"

	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/dm"
	"github.com/weaveworks/ignite/pkg/format"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/source"
	"github.com/weaveworks/ignite/pkg/util"
)

const (
	dataDevSize = 100 * 1073741824 // 100 GB
	blockSize   = 128              // Pool allocation block size
)

// Additional space to add to the image for the ext4 partition (100 MB)
var extraSize = format.SectorsFromBytes(100 * 1048576)

func (md *ImageMetadata) NewDMPool() error {
	metadataFile := path.Join(md.ObjectPath(), constants.IMAGE_THINMETADATA)
	dataFile := path.Join(md.ObjectPath(), constants.IMAGE_THINDATA)

	if err := allocateDMFiles(metadataFile, dataFile); err != nil {
		return err
	}

	md.ImageOD().Pool = dm.NewPool(
		util.NewPrefixer().Prefix("pool", md.ID.String()),
		format.SectorsFromBytes(dataDevSize),
		format.Sectors(blockSize),
		dm.NewLoopDevice(metadataFile, false),
		dm.NewLoopDevice(dataFile, false),
	)

	return nil
}

// Allocate the thin provisioning data and metadata files
func allocateDMFiles(metadataFile, dataFile string) error {
	thinFiles := map[string]int64{
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
			if err := file.Truncate(size); err != nil {
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

	volume, err := od.Pool.CreateVolume(p.Prefix(src.ID()), src.SizeSectors()+extraSize)
	if err != nil {
		return nil, err
	}

	mountPoint, err := volume.Import(src)
	if err != nil {
		return nil, err
	}

	//tarCmd := exec.Command("tar", "-x", "-C", mountPoint.Path)
	//reader, err := src.Reader()
	//if err != nil {
	//	return nil, err
	//}
	//
	//tarCmd.Stdin = reader
	//if err := tarCmd.Start(); err != nil {
	//	return nil, err
	//}
	//
	//if err := tarCmd.Wait(); err != nil {
	//	return nil, err
	//}
	//
	//if err := src.Cleanup(); err != nil {
	//	return nil, err
	//}

	// Test kernel import
	//if err := testKernelImport(volume, p); err != nil {
	//	return nil, err
	//}

	// Test snapshot creation
	//if _, err := volume.CreateSnapshot("snapshot-" + src.ID(), volume.Size()); err != nil {
	//	return nil, err
	//}
	//
	//if _, err := volume.CreateSnapshot("snapshot2-" + src.ID(), volume.Size()); err != nil {
	//	return nil, err
	//}
	//
	//od.Pool.Remove(1)

	return mountPoint, nil
}

func (md *ImageMetadata) CreateOverlay(kernelSrc source.Source, size uint64, id *metadata.ID) (*dm.Device, error) {
	sizeSectors := format.SectorsFromBytes(size)
	p := util.NewPrefixer()
	volumeName := p.Prefix(md.ID.String())
	resizeName := p.Prefix("resize", sizeSectors.String())
	kernelName := p.Prefix("kernel")
	overlayName := p.Prefix(id.String())

	// Requested kernel doesn't exist, so import it
	var kernel *dm.Device
	var err error
	if kernel, err = md.ImageOD().Pool.Get(kernelName); err != nil {
		volume, err := md.ImageOD().Pool.Get(volumeName)
		if err != nil {
			return nil, err
		}

		resize, err := volume.CreateSnapshot(resizeName, sizeSectors)
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

func testKernelImport(volume *dm.Device, p *util.Prefixer) error {
	// TODO: These are the test case variables
	var resizeSize = format.SectorsFromBytes(10 * 1073741824)

	resize, err := volume.CreateSnapshot(p.Prefix("resize", resizeSize.String()), resizeSize)
	if err != nil {
		return err
	}

	// The kernel
	kernel, err := resize.CreateSnapshot(p.Prefix("kernel"), resize.Size())
	if err != nil {
		return err
	}

	kernelSrc, err := source.NewDockerSource("weaveworks/ignite-kernel:4.19.47")
	if err != nil {
		return err
	}

	mountPoint, err := kernel.Import(kernelSrc)
	if err != nil {
		return err
	}

	if err := mountPoint.Umount(); err != nil {
		return err
	}

	return nil
}

func calcMetadataDevSize(dataDevSize int64) int64 {
	// The minimum size is 2 MB and the maximum size is 16 GB
	var minSize int64 = 2 * 1048576
	var maxSize int64 = 16 * 1073741824
	size := 48 * dataDevSize / blockSize

	if size < minSize {
		return minSize
	} else if size > maxSize {
		return maxSize
	}

	return size
}
