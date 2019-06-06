package imgmd

import (
	"archive/tar"
	"io/ioutil"
	"path/filepath"
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/util"
	"github.com/mholt/archiver"
	"github.com/nerd2/gexto"
	"github.com/pkg/errors"
	"os"
	"path"
)

func (md *ImageMetadata) ImportImage(p string) error {
	if err := util.CopyFile(p, path.Join(md.ObjectPath(), constants.IMAGE_FS)); err != nil {
		return fmt.Errorf("failed to copy image file %q to image %q: %v", p, md.ID, err)
	}

	return nil
}

func (md *ImageMetadata) AllocateAndFormat(size int64) error {
	p := path.Join(md.ObjectPath(), constants.IMAGE_FS)
	imageFile, err := os.Create(p)
	if err != nil {
		return errors.Wrapf(err, "failed to create image file for %s", md.ID)
	}
	defer imageFile.Close()

	// The base image is the size of the tar file, plus 100MB
	// TODO: This is temporary only, until we make DM snapshot overlays work.
	if err := imageFile.Truncate(size + 4000 * 1024 * 1024); err != nil {
		return errors.Wrapf(err, "failed to allocate space for image %s", md.ID)
	}

	// Use mkfs.ext4 to create the new image with an inode size of 256
	// (gexto doesn't support anything but 128, but as long as we're not using that it's fine)
	if _, err := util.ExecuteCommand("mkfs.ext4", "-I", "256", "-E", "lazy_itable_init=0,lazy_journal_init=0", p); err != nil {
		return errors.Wrapf(err, "failed to format image %s", md.ID)
	}

	return nil
}

// AddFilesWithGexto adds all the files from the given rootfs tar to the image
// TODO: Fix the "corrupt direntry" error from gexto
func (md *ImageMetadata) AddFilesWithGexto(sourcePath string) error {
	// TODO: This
	p := path.Join(md.ObjectPath(), constants.IMAGE_FS)
	filesystem, err := gexto.NewFileSystem(p)
	if err != nil {
		return err
	}
	//defer filesystem.Close()

	if err := archiver.Walk(sourcePath, func(f archiver.File) error {
		th, ok := f.Header.(*tar.Header)

		if !ok {
			return fmt.Errorf("expected header to be *tar.Header but was %T", f.Header)
		}

		//f.FileInfo
		//data, err := os.Stat(f.Name())
		//if err != nil {
		//	return err
		//}

		relativePath, _ := filepath.Rel(".", th.Name)
		filePath := filepath.Join("/", relativePath)

		if f.IsDir() {
			fmt.Printf("Directory: %s\n", filePath)
			if err := filesystem.Mkdir(filePath, 0777); err != nil {
				return fmt.Errorf("unable to create directory in filesystem: %s", filePath)
			}
		} else {
			file, err := filesystem.Create(filePath)
			if err != nil {
				return fmt.Errorf("unable to create file in filesystem: %s", filePath)
			}
			//
			//contents, err := ioutil.ReadAll(f)
			//if err != nil {
			//	return fmt.Errorf("unable to read file contents: %s", filePath)
			//}
			file.Write([]byte("Hello, world!"))
			f.Close()
		}

		//if ok {
		//	fmt.Println("Filename:", zfh.Name)
		//}

		//if f.IsDir() {
		//	fmt.Println("Directory:", f.Name())
		//} else {
		//	fmt.Println("Filename:", f.Name())
		//}
		return nil
	}); err != nil {
		return err
	}

	filesystem.Close()
	return nil
}

// AddFiles copies the contents of the tar file into the ext4 filesystem
func (md *ImageMetadata) AddFiles(sourcePath string) error {
	p := path.Join(md.ObjectPath(), constants.IMAGE_FS)
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	if _, err := util.ExecuteCommand("mount", "-o", "loop", p, tempDir); err != nil {
		return fmt.Errorf("failed to mount image %q: %v", p, err)
	}
	defer util.ExecuteCommand("umount", tempDir)

	if _, err := util.ExecuteCommand("tar", "-xf", sourcePath, "-C", tempDir); err != nil {
		return err
	}
	return md.SetupResolvConf(tempDir)
}

// SetupResolvConf makes sure there is a resolv.conf file, otherwise
// name resolution won't work. The kernel uses DHCP by default, and
// puts the nameservers in /proc/net/pnp at runtime. Hence, as a default,
// if /etc/resolv.conf doesn't exist, we can use /proc/net/pnp as /etc/resolv.conf
func (md *ImageMetadata) SetupResolvConf(tempDir string) error {
	resolvConf := filepath.Join(tempDir, "/etc/resolv.conf")
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

func (md *ImageMetadata) ExportKernel() (string, error) {
	p := path.Join(md.ObjectPath(), constants.IMAGE_FS)
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}

	kernelDir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}

	if _, err := util.ExecuteCommand("mount", "-o", "loop", p, tempDir); err != nil {
		return "", fmt.Errorf("failed to mount image %q: %v", p, err)
	}
	defer util.ExecuteCommand("umount", tempDir)

	kernelDest := path.Join(kernelDir, constants.KERNEL_FILE)
	kernelSrc, err := findKernel(tempDir)
	if err != nil {
		return "", err
	}

	if util.FileExists(kernelSrc) {
		if err := util.CopyFile(kernelSrc, kernelDest); err != nil {
			return "", fmt.Errorf("failed to copy kernel file from %q to %q: %v", kernelSrc, kernelDest, err)
		}
	} else {
		return "", fmt.Errorf("no kernel found in image %q", md.ID)
	}

	return kernelDir, nil
}

func (md *ImageMetadata) Size() (int64, error) {
	fi, err := os.Stat(path.Join(md.ObjectPath(), constants.IMAGE_FS))
	if err != nil {
		return 0, err
	}

	return fi.Size(), nil
}

// Quick hack to resolve a kernel in the image
func findKernel(tmpDir string) (string, error) {
	bootDir := path.Join(tmpDir, "boot")
	kernel := path.Join(bootDir, constants.KERNEL_FILE)

	fi, err := os.Lstat(kernel)
	if err != nil {
		return "", err
	}

	// The target is a symlink
	if fi.Mode()&os.ModeSymlink != 0 {
		kernel, err = os.Readlink(kernel)
		if err != nil {
			return "", err
		}

		// Fix the path for absolute and relative symlinks
		if path.IsAbs(kernel) {
			kernel = path.Join(tmpDir, kernel)
		} else {
			kernel = path.Join(bootDir, kernel)
		}
	}

	return kernel, nil
}
