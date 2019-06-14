package imgmd

import (
	"os"
	"path"
	"path/filepath"

	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/source"
	"github.com/weaveworks/ignite/pkg/util"
)

func (md *ImageMetadata) CreateImageFile(size int64) (*source.ImageFile, error) {
	// Add 100 MB to the tar file size to be safe
	return source.NewImageFile(path.Join(md.ObjectPath(), constants.IMAGE_FS), size+100*1048576)
}

// AddFiles copies the contents of the tar file into the ext4 filesystem
func (md *ImageMetadata) AddFiles(i *source.ImageFile, src source.Source) error {
	mountPoint, err := i.AddFiles(src)
	if err != nil {
		return err
	}
	defer mountPoint.Umount()

	// Check if this is a "combined image" now that it's mounted
	if kernel, err := util.FindKernel(mountPoint.Path); err == nil {
		md.ImageOD().ContainsKernel = len(kernel) > 0
	} else {
		return err
	}

	return md.setupResolvConf(mountPoint.Path)
}

// setupResolvConf makes sure there is a resolv.conf file, otherwise
// name resolution won't work. The kernel uses DHCP by default, and
// puts the nameservers in /proc/net/pnp at runtime. Hence, as a default,
// if /etc/resolv.conf doesn't exist, we can use /proc/net/pnp as /etc/resolv.conf
func (md *ImageMetadata) setupResolvConf(dir string) error {
	resolvConf := filepath.Join(dir, "/etc/resolv.conf")
	empty, err := util.FileIsEmpty(resolvConf)
	if err != nil {
		return err
	}

	if !empty {
		return nil
	}

	//fmt.Println("Symlinking /etc/resolv.conf to /proc/net/pnp")
	return os.Symlink("../proc/net/pnp", resolvConf)
}

func (md *ImageMetadata) Size() (int64, error) {
	imageFile, err := source.LoadImageFile(path.Join(md.ObjectPath(), constants.IMAGE_FS))
	if err != nil {
		return 0, err
	}

	return imageFile.Size(), nil
}
