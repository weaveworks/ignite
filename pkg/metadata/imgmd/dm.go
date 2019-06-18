package imgmd

import (
	"fmt"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/dm"
	"github.com/weaveworks/ignite/pkg/source"
	"github.com/weaveworks/ignite/pkg/util"
	"os"
	"os/exec"
	"path"
)

const (
	dataDevSize = 100 * 1073741824 // 100 GB
	blockSize   = 128              // Pool allocation block size
)

var imageExtraSize = dm.SectorsFromBytes(100 * 1048576)

func prefix(input string) string {
	return constants.IGNITE_PREFIX + input
}

type ImageDM struct {
	*dm.DMPool
}

func emptyImageDM() *ImageDM {
	return &ImageDM{
		&dm.DMPool{},
	}
}

func (md *ImageMetadata) NewImageDM() error {
	metadataFile := path.Join(md.ObjectPath(), constants.IMAGE_THINMETADATA)
	dataFile := path.Join(md.ObjectPath(), constants.IMAGE_THINDATA)

	if err := allocateDMFiles(metadataFile, dataFile); err != nil {
		return err
	}

	md.ImageOD().ImageDM = &ImageDM{dm.NewDMPool(
		prefix("pool-"+md.ID.String()),
		dm.SectorsFromBytes(dataDevSize),
		dm.Sectors(blockSize),
		dm.NewLoopDevice(metadataFile, false),
		dm.NewLoopDevice(dataFile, false),
	)}

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

func (dm *ImageDM) AddFiles(src source.Source) (*util.MountPoint, error) {
	volume, err := dm.CreateVolume(src.ID(), src.SizeSectors()+imageExtraSize)
	if err != nil {
		return nil, err
	}

	if _, err := util.ExecuteCommand("mkfs.ext4", "-I", "256",
		"-E", "lazy_itable_init=0,lazy_journal_init=0", volume.Path()); err != nil {
		return nil, fmt.Errorf("failed to format image %q: %v", volume.Path(), err)
	}

	mountPoint, err := util.Mount(volume.Path())
	if err != nil {
		return nil, err
	}

	tarCmd := exec.Command("tar", "-x", "-C", mountPoint.Path)
	reader, err := src.Reader()
	if err != nil {
		return nil, err
	}

	tarCmd.Stdin = reader
	if err := tarCmd.Start(); err != nil {
		return nil, err
	}

	if err := tarCmd.Wait(); err != nil {
		return nil, err
	}

	if err := src.Cleanup(); err != nil {
		return nil, err
	}

	// Test snapshot creation
	if _, err := volume.CreateSnapshot("snapshot-" + src.ID()); err != nil {
		return nil, err
	}

	if _, err := volume.CreateSnapshot("snapshot2-" + src.ID()); err != nil {
		return nil, err
	}

	dm.Remove(1)

	return mountPoint, nil
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
