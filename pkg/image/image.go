package build

import (
	"archive/tar"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

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

type ImageMetadata struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
}

func (i ImageMetadata) AllocateAndFormat() error {
	p := path.Join(constants.IMAGE_DIR, i.ID, constants.IMAGE_FS)
	imageFile, err := os.Create(p)
	if err != nil {
		return errors.Wrapf(err, "failed to create image file for %s", i.ID)
	}
	defer imageFile.Close()

	// TODO: Dynamic size, for now hardcoded 4 GiB
	if err := imageFile.Truncate(4294967296); err != nil {
		return errors.Wrapf(err, "failed to allocate space for image %s", i.ID)
	}

	//blank := make([]byte, 1024*1024)
	//for i := 0; i < 4096; i++ {
	//	_, _ = imageFile.Write(blank)
	//}
	//_ = imageFile.Close()

	// Use mkfs.ext4 to create the new image with an inode size of 128 (gexto doesn't support the default of 256)
	if _, err := util.ExecuteCommand("mkfs.ext4", "-I", "128", "-E", "lazy_itable_init=0,lazy_journal_init=0", p); err != nil {
		return errors.Wrapf(err, "failed to format image %s", i.ID)
	}

	return nil
}

// Adds all the files from the given rootfs tar to the image
// TODO: Fix the "corrupt direntry" error from gexto
func (i ImageMetadata) AddFiles(sourcePath string) error {
	// TODO: This
	p := path.Join(constants.IMAGE_DIR, i.ID, constants.IMAGE_FS)
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
func (i ImageMetadata) AddFiles2(sourcePath string) error {
	p := path.Join(constants.IMAGE_DIR, i.ID, constants.IMAGE_FS)
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

func (i ImageMetadata) WriteMetadata() error {
	f, err := os.Create(path.Join(constants.IMAGE_DIR, i.ID, constants.METADATA))
	if err != nil {
		return err
	}
	defer f.Close()

	y, err := json.MarshalIndent(&i, "", "    ")
	if err != nil {
		return err
	}

	if _, err := f.Write(append(y, '\n')); err != nil {
		return err
	}

	return nil
}
