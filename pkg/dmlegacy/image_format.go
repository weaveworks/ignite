package dmlegacy

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/source"
	"github.com/weaveworks/ignite/pkg/util"
)

const blockSize = 4096 // Block size to use for the ext4 filesystems, this is the default

// CreateImageFilesystem creates an ext4 filesystem in a file, containing the files from the source
func CreateImageFilesystem(img *api.Image, src source.Source) error {
	log.Debugf("Allocating image file and formatting it with ext4...")
	p := path.Join(img.ObjectPath(), constants.IMAGE_FS)
	imageFile, err := os.Create(p)
	if err != nil {
		return errors.Wrapf(err, "failed to create image file for %s", img.GetUID())
	}
	defer imageFile.Close()

	// To accommodate space for the tar file contents and the ext4 journal + other metadata,
	// make the base image a sparse file three times the size of the source contents. This
	// will be shrunk to fit by resizeToMinimum later.
	if err := imageFile.Truncate(int64(img.Status.OCISource.Size.Bytes()) * 3); err != nil {
		return errors.Wrapf(err, "failed to allocate space for image %s", img.GetUID())
	}

	// Use mkfs.ext4 to create the new image with an inode size of 256
	// (gexto doesn't support anything but 128, but as long as we're not using that it's fine)
	if _, err := util.ExecuteCommand("mkfs.ext4", "-b", strconv.Itoa(blockSize),
		"-I", "256", "-F", "-E", "lazy_itable_init=0,lazy_journal_init=0", p); err != nil {
		return errors.Wrapf(err, "failed to format image %s", img.GetUID())
	}

	// Proceed with populating the image with files
	if err := addFiles(img, src); err != nil {
		return err
	}

	// Resize the image to its minimum size
	return resizeToMinimum(img)
}

// addFiles copies the contents of the tar file into the ext4 filesystem
func addFiles(img *api.Image, src source.Source) error {
	log.Debugf("Copying in files to the image file from a source...")
	p := path.Join(img.ObjectPath(), constants.IMAGE_FS)
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

	return setupResolvConf(tempDir)
}

// setupResolvConf makes sure there is a resolv.conf file, otherwise
// name resolution won't work. The kernel uses DHCP by default, and
// puts the nameservers in /proc/net/pnp at runtime. Hence, as a default,
// if /etc/resolv.conf doesn't exist, we can use /proc/net/pnp as /etc/resolv.conf
func setupResolvConf(tempDir string) error {
	resolvConf := filepath.Join(tempDir, "/etc/resolv.conf")
	empty, err := util.FileIsEmpty(resolvConf)
	if err != nil {
		return err
	}

	if !empty {
		return nil
	}

	return os.Symlink("../proc/net/pnp", resolvConf)
}

// resizeToMinimum resizes the given image to the smallest size possible
func resizeToMinimum(img *api.Image) (err error) {
	p := path.Join(img.ObjectPath(), constants.IMAGE_FS)
	var minSize int64
	var imageFile *os.File

	if minSize, err = getMinSize(p); err != nil {
		return
	}

	if imageFile, err = os.OpenFile(p, os.O_RDWR, constants.DATA_DIR_FILE_PERM); err != nil {
		return err
	}
	defer util.DeferErr(&err, imageFile.Close)

	minSizeBytes := minSize * blockSize

	log.Debugf("Truncating %q to %d bytes", p, minSizeBytes)
	if err = imageFile.Truncate(minSizeBytes); err != nil {
		err = fmt.Errorf("failed to shrink image %q: %v", img.GetUID(), err)
	}

	return
}

// getMinSize retrieves the minimum size for a block device file
// containing a filesystem and shrinks the filesystem to that size
func getMinSize(p string) (minSize int64, err error) {
	const (
		Colon          = ":"
		FullWidthColon = "："
	)

	// Loop mount the image for resize2fs
	imageLoop, err := newLoopDev(p, false)
	if err != nil {
		return
	}

	// Defer the detach
	defer util.DeferErr(&err, imageLoop.Detach)

	// Call e2fsck for resize2fs, it sometimes requires this
	// e2fsck throws an error if the filesystem gets repaired, so just ignore it
	_, _ = util.ExecuteCommand("e2fsck", "-p", "-f", imageLoop.Path())

	// Retrieve the minimum size for the filesystem
	log.Debugf("Retrieving minimum size for %q", imageLoop.Path())
	out, err := util.ExecuteCommand("resize2fs", "-P", imageLoop.Path())
	if err != nil {
		return
	}

	// resize2fs 1.45.3 (14-Jul-2019)
	// Estimated minimum size of the filesystem: 5813528
	// Try split with colon
	split := strings.SplitN(out, Colon, 2)
	if len(split) != 2 {
		// Seems to fail, try split again with full width colon
		split = strings.SplitN(out, FullWidthColon, 2)
		if len(split) != 2 {
			err = errors.Errorf("cannot parse minimum size from resize2fs: %s", out)
			return
		}
	}
	minSizeToParse := strings.TrimSpace(split[len(split)-1])
	if minSize, err = strconv.ParseInt(minSizeToParse, 10, 64); err != nil {
		return
	}

	log.Debugf("Minimum size: %d blocks", minSize)

	// Perform the filesystem resize
	_, err = util.ExecuteCommand("resize2fs", imageLoop.Path(), strconv.FormatInt(minSize, 10))
	return
}
