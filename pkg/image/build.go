package build

import (
	"crypto/rand"
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/util"
	"github.com/pkg/errors"
	"path"
)

//import (
//	"os"
//	"path/filepath"
//	"strings"
//
//	"github.com/pkg/errors"
//)

type VMID []byte

/*
func (vmm *VMM) copyFilesFromHost() error {
	if len(vmm.copyFiles) == 0 {
		return nil
	}
	mntdir := filepath.Join(RuntimeDir, vmm.name, "mnt")
	if err := os.MkdirAll(mntdir, 0755); err != nil {
		return err
	}
	if err := executeCommand("mount", vmm.rootDrivePath, mntdir); err != nil {
		return err
	}
	for _, filePair := range vmm.copyFiles {
		files := strings.Split(filePair, ":")
		if len(files) != 2 {
			return errors.Errorf("--copy-files arguments must be of the form SOURCE:TARGET")
		}
		src := files[0]
		dest := filepath.Join(mntdir, files[1])
		destDir := filepath.Dir(dest)
		if _, err := os.Stat(destDir); os.IsNotExist(err) {
			if err := os.MkdirAll(destDir, 755); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		if err := executeCommand("cp", "-r", src, dest); err != nil {
			return err
		}
	}
	if err := executeCommand("sync", vmm.rootDrivePath); err != nil {
		return err
	}
	return executeCommand("umount", mntdir)
}
*/

//func NewImage(vmID string) {
//	image := path.Join(constants.VM_DIR, vmID, constants.VM_FS_IMAGE)
//
//}

// Creates a new 8-byte VM ID and return it as a string
func NewVMID() (string, error) {
	var vmID string
	var idBytes []byte

	for {
		idBytes = make([]byte, 8)
		if _, err := rand.Read(idBytes); err != nil {
			return "", errors.Wrap(err, "failed to generate VM ID")
		}

		// Convert the byte array to a string literally
		vmID = fmt.Sprintf("%x", idBytes)

		// If the generated ID is unique, return it
		if exists, _ := util.PathExists(path.Join(constants.VM_DIR, vmID)); !exists {
			return vmID, nil
		}
	}
}
