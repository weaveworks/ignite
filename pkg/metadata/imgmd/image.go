package imgmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/source"
	"github.com/weaveworks/ignite/pkg/util"
)

func (md *Image) AllocateAndFormat() error {
	log.Debugf("Allocating image file and formatting it with ext4...")
	p := path.Join(md.ObjectPath(), constants.IMAGE_FS)
	imageFile, err := os.Create(p)
	if err != nil {
		return errors.Wrapf(err, "failed to create image file for %s", md.GetUID())
	}
	defer imageFile.Close()

	// The base image is the size of the tar file, plus 100MB
	if err := imageFile.Truncate(int64(md.Status.OCISource.Size.Bytes()) + 100*1024*1024); err != nil {
		return errors.Wrapf(err, "failed to allocate space for image %s", md.GetUID())
	}

	// Use mkfs.ext4 to create the new image with an inode size of 256
	// (gexto doesn't support anything but 128, but as long as we're not using that it's fine)
	if _, err := util.ExecuteCommand("mkfs.ext4", "-I", "256", "-F", "-E", "lazy_itable_init=0,lazy_journal_init=0", p); err != nil {
		return errors.Wrapf(err, "failed to format image %s", md.GetUID())
	}

	return nil
}

// AddFiles copies the contents of the tar file into the ext4 filesystem
func (md *Image) AddFiles(src source.Source) error {
	log.Debugf("Copying in files to the image file from a source...")
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

	tarCmd := exec.Command("tar", "-x", "-C", tempDir)
	reader, err := src.Reader()
	if err != nil {
		return err
	}

	tarCmd.Stdin = reader
	if err := tarCmd.Start(); err != nil {
		return err
	}

	if err := tarCmd.Wait(); err != nil {
		return err
	}

	if err := src.Cleanup(); err != nil {
		return err
	}

	return md.SetupResolvConf(tempDir)
}

// SetupResolvConf makes sure there is a resolv.conf file, otherwise
// name resolution won't work. The kernel uses DHCP by default, and
// puts the nameservers in /proc/net/pnp at runtime. Hence, as a default,
// if /etc/resolv.conf doesn't exist, we can use /proc/net/pnp as /etc/resolv.conf
func (md *Image) SetupResolvConf(tempDir string) error {
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
