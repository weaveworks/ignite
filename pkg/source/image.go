package source

import (
	"fmt"
	"github.com/weaveworks/ignite/pkg/dm"
	"os"
	"os/exec"

	"github.com/weaveworks/ignite/pkg/util"
)

type ImageFile struct {
	path string
	size int64
}

func NewImageFile(path string, size int64) (*ImageFile, error) {
	i := &ImageFile{
		path: path,
		size: size,
	}

	if err := i.Create(); err != nil {
		return nil, err
	}

	return i, nil
}

func LoadImageFile(path string) (*ImageFile, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	return &ImageFile{
		path: path,
		size: fi.Size(),
	}, nil
}

// Get the image size as bytes
func (i *ImageFile) SizeBytes() int64 {
	return i.size
}

// Get the image size as 512-byte sectors
func (i *ImageFile) SizeSectors() dm.Sectors {
	return dm.SectorsFromBytes(i.size)
}

// Create allocates and formats the image file
func (i *ImageFile) Create() error {
	file, err := os.Create(i.path)
	if err != nil {
		return fmt.Errorf("failed to create image file %q: %v", i.path, err)
	}
	defer file.Close()

	// Align the image to 512-byte sectors, otherwise loop mounting doesn't work
	//size = int64(math.Ceil(float64(size)/512) * 512)

	// Allocate the image file
	if err := file.Truncate(i.size); err != nil {
		return fmt.Errorf("failed to allocate space for image %q: %v", i.path, err)
	}

	// Use mkfs.ext4 to create the new image file with an inode size of 256
	// (gexto supports only 128, but as long as we're not using it for the moment)
	if _, err := util.ExecuteCommand("mkfs.ext4", "-I", "256",
		"-E", "lazy_itable_init=0,lazy_journal_init=0", i.path); err != nil {
		return fmt.Errorf("failed to format image %q: %v", i.path, err)
	}

	return nil
}

func (i *ImageFile) AddFiles(src Source) (*util.MountPoint, error) {
	mountPoint, err := util.Mount(i.path)
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

	return mountPoint, nil
}
