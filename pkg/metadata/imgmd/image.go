package imgmd

import (
	"archive/tar"
	"io/ioutil"
	"path/filepath"
	"strings"
	//"archive/tar"
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

func (md *ImageMetadata) AllocateAndFormat() error {
	p := path.Join(md.ObjectPath(), constants.IMAGE_FS)
	imageFile, err := os.Create(p)
	if err != nil {
		return errors.Wrapf(err, "failed to create image file for %s", md.ID)
	}
	defer imageFile.Close()

	// TODO: Dynamic size, for now hardcoded 4 GiB
	if err := imageFile.Truncate(4294967296); err != nil {
		return errors.Wrapf(err, "failed to allocate space for image %s", md.ID)
	}

	//blank := make([]byte, 1024*1024)
	//for i := 0; i < 4096; i++ {
	//	_, _ = imageFile.Write(blank)
	//}
	//_ = imageFile.Close()

	// Use mkfs.ext4 to create the new image with an inode size of 128 (gexto doesn't support the default of 256)
	if _, err := util.ExecuteCommand("mkfs.ext4", "-I", "128", "-E", "lazy_itable_init=0,lazy_journal_init=0", p); err != nil {
		return errors.Wrapf(err, "failed to format image %s", md.ID)
	}

	return nil
}

// Adds all the files from the given rootfs tar to the image
// TODO: Fix the "corrupt direntry" error from gexto
func (md *ImageMetadata) AddFiles(sourcePath string) error {
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

	fmt.Println("Open success!")

	return nil
	// We need to open each file in the rootfs.tar and write it into fs
}

// mount-based file adder (temporary, requires root)
func (md *ImageMetadata) AddFiles2(sourcePath string) error {
	p := path.Join(md.ObjectPath(), constants.IMAGE_FS)
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	if _, err := util.ExecuteCommand("sudo", "mount", "-o", "loop", p, tempDir); err != nil {
		return errors.Wrapf(err, "failed to mount image %s", p)
	}
	defer util.ExecuteCommand("sudo", "umount", tempDir)

	if err := archiver.Unarchive(sourcePath, tempDir); err != nil {
		return err
	}

	return nil
}

func (md *ImageMetadata) AddFiles3(sourcePath string) error {
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

	return nil
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
	kernelSrc, err := findKernel(path.Join(tempDir, "boot"))
	if err != nil {
		return "", err
	}

	if kernelSrc != "" {
		if err := util.CopyFile(kernelSrc, kernelDest); err != nil {
			return "", fmt.Errorf("failed to copy kernel file from %q to %q: %v", kernelSrc, kernelDest, err)
		}
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
func findKernel(dir string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "vmlinux-") {
			kernel := path.Join(dir, file.Name())
			fmt.Printf("Found a kernel: %q\n", file.Name())
			return kernel, nil
		}
	}

	// TODO: Replace these with logging
	fmt.Println("No kernel found in image")

	return "", nil
}
