package kernmd

import (
	"fmt"
	"os"
	"path"

	"github.com/weaveworks/ignite/pkg/source"

	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/util"
)

func (md *KernelMetadata) ImportKernel(p string) error {
	if err := util.CopyFile(p, path.Join(md.ObjectPath(), constants.KERNEL_FILE)); err != nil {
		return fmt.Errorf("failed to copy kernel file %q to kernel %q: %v", p, md.ID, err)
	}

	return nil
}

func (md *KernelMetadata) CreateImageFile(size int64) (*source.ImageFile, error) {
	// Add 100 MB to the tar file size to be safe
	return source.NewImageFile(path.Join(md.ObjectPath(), constants.KERNEL_FS), size+100*1048576)
}

// AddFiles copies the contents of the tar file into the ext4 filesystem
func (md *KernelMetadata) AddFiles(i *source.ImageFile, src source.Source) error {
	mountPoint, err := i.AddFiles(src)
	if err != nil {
		return err
	}
	defer mountPoint.Umount()

	return md.exportVMLinux(mountPoint.Path)
}

// Copies vmlinux out of the kernel image for Firecracker
func (md *KernelMetadata) exportVMLinux(dir string) error {
	kernelFile, err := util.FindKernel(dir)
	if err != nil {
		return err
	}

	return util.CopyFile(kernelFile, path.Join(md.ObjectPath(), constants.KERNEL_FILE))
}

// Gets the size of the kernel filesystem
func (md *KernelMetadata) Size() (int64, error) {
	fi, err := os.Stat(path.Join(md.ObjectPath(), constants.KERNEL_FS))
	if err != nil {
		return 0, err
	}

	return fi.Size(), nil
}
