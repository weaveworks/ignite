package util

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"io"
	"io/ioutil"
	"os"
)

// Creates the /var/lib/firecracker/{vm,image,kernel} directories
func CreateDirectories() error {
	for _, dir := range []string{constants.VM_DIR, constants.IMAGE_DIR, constants.KERNEL_DIR} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %q: %v", dir, err)
		}
	}

	return nil
}

func PathExists(path string) (bool, os.FileInfo) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, info
}

func FileExists(filename string) bool {
	exists, info := PathExists(filename)
	if !exists {
		return false
	}
	return !info.IsDir()
}

func DirExists(dirname string) bool {
	exists, info := PathExists(dirname)
	if !exists {
		return false
	}
	return info.IsDir()
}

func CopyFile(src string, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

type MountPoint struct {
	Path string
}

func Mount(volume string) (*MountPoint, error) {
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}

	if _, err := ExecuteCommand("mount", volume, tempDir); err != nil {
		return nil, fmt.Errorf("failed to mount volume %q: %v", volume, err)
	}

	return &MountPoint{
		Path: tempDir,
	}, nil
}

func (mp *MountPoint) Umount() error {
	if _, err := ExecuteCommand("umount", mp.Path); err != nil {
		return fmt.Errorf("failed to unmount volume %q: %v", mp.Path, err)
	}

	if err := os.RemoveAll(mp.Path); err != nil {
		return err
	}

	return nil
}
