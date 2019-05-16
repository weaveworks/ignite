package build

//import (
//	"os"
//	"path/filepath"
//	"strings"
//
//	"github.com/pkg/errors"
//)

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
