package util

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/otiai10/copy"
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/constants"
)

// Creates the /var/lib/firecracker/{vm,image,kernel} directories
func CreateDirectories() error {
	for _, dir := range []string{constants.VM_DIR, constants.IMAGE_DIR, constants.KERNEL_DIR, constants.MANIFEST_DIR} {
		if err := os.MkdirAll(dir, constants.DATA_DIR_PERM); err != nil {
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

func DirEmpty(dirname string) (b bool) {
	if !DirExists(dirname) {
		return
	}

	f, err := os.Open(dirname)
	if err != nil {
		return
	}
	defer func() { _ = f.Close() }()

	// If the first file is EOF, the directory is empty
	if _, err = f.Readdir(1); err == io.EOF {
		b = true
	}

	return
}

func IsDeviceFile(filename string) (err error) {
	if exists, info := PathExists(filename); !exists {
		err = fmt.Errorf("device path not found")
	} else if info.Mode()&os.ModeDevice == 0 {
		err = fmt.Errorf("not a device file")
	}

	return
}

// CopyFile copies both files and directories
func CopyFile(src string, dst string) error {
	return copy.Copy(src, dst)
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

// FileIsEmpty returns true if the file is empty
func FileIsEmpty(file string) (bool, error) {
	fileInfo, err := os.Stat(file)
	// Check if there was an unexpected error
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}

	// The file exists, and has content. Proceed as usual
	if err == nil && fileInfo.Size() > 0 {
		return false, nil
	}

	// The file exists, but has no content. Remove the file to allow the symlink
	if err == nil && fileInfo.Size() == 0 {
		if err := os.Remove(file); err != nil {
			return false, err
		}
	}

	return true, nil
}

// WriteFileIfChanged stores a sha of data at <filename>.sha256 and determines whether to
// rewrite the file; it has the same signature as ioutil.WriteFile().
func WriteFileIfChanged(filename string, data []byte, perm os.FileMode) error {
	shaFile := filename + ".sha256"

	currentHashBytes, err := ioutil.ReadFile(shaFile)
	currentHashStr := string(currentHashBytes)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	newHashHex := sha256.Sum256(data)
	newHashStr := hex.EncodeToString(newHashHex[:])

	if newHashStr != currentHashStr {
		log.Debugf("Writing %q with new hash: %q, old hash: %q", filename, newHashStr, currentHashStr)
		err = ioutil.WriteFile(
			shaFile,
			[]byte(newHashStr),
			perm,
		)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(
			filename,
			data,
			perm,
		)
	}
	log.Debugf("%q with hash %q is unchanged", filename, newHashStr)

	return nil
}
