package imgmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/util"
)

type ImageSource struct {
	size        int64
	dockerImage string
	imageID     string
	containerID string
	dockerCmd   *exec.Cmd
	tarFile     string
}

func NewSource(src string) (*ImageSource, error) {
	if util.FileExists(src) {
		if !strings.HasSuffix(src, ".tar") {
			// TODO: Allow reading from stdin
			return nil, fmt.Errorf("only tar files or docker images are supported import methods")
		}
		fo, err := os.Stat(src)
		if err != nil {
			return nil, err
		}
		return &ImageSource{
			tarFile: src,
			size:    fo.Size(),
		}, nil
	}
	// Treat the source as a docker image
	// If it doesn't have a tag, assume it's latest
	if !strings.Contains(src, ":") {
		src += ":latest"
	}
	// Query docker for the image
	out, err := util.ExecuteCommand("docker", "images", "-q", src)
	if err != nil {
		return nil, err
	}
	// TODO: docker pull if it's not found
	if util.IsEmptyString(out) {
		return nil, fmt.Errorf("docker image %s not found", src)
	}

	// Docker outputs one image per line
	dockerIDs := strings.Split(strings.TrimSpace(out), "\n")

	// Check if the image query is too broad
	if len(dockerIDs) > 1 {
		return nil, fmt.Errorf("multiple matches, too broad docker image query: %s", src)
	}

	// Select the first (and only) match
	dockerID := dockerIDs[0]

	// docker inspect quay.io/footloose/centos7ignite -f "{{.Size}}"
	out, err = util.ExecuteCommand("docker", "inspect", src, "-f", "{{.Size}}")
	if err != nil {
		return nil, err
	}
	size, err := strconv.Atoi(out)
	if err != nil {
		return nil, err
	}
	return &ImageSource{
		dockerImage: src,
		imageID:     dockerID,
		size:        int64(size),
	}, nil
}

func (is *ImageSource) GetReader() (io.ReadCloser, error) {
	if len(is.imageID) > 0 {
		// Create a container from the image to be exported
		containerID, err := util.ExecuteCommand("docker", "create", is.imageID, "sh")
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create docker container from image %s", is.imageID)
		}

		// Export the created container to a tar archive that will be later extracted into the VM disk image
		is.dockerCmd = exec.Command("docker", "export", containerID)
		reader, err := is.dockerCmd.StdoutPipe()
		if err != nil {
			return nil, err
		}
		if err := is.dockerCmd.Start(); err != nil {
			return nil, err
		}
		return reader, nil
	}
	return os.Open(is.tarFile)
}

func (is *ImageSource) DockerImage() string {
	return is.dockerImage
}

func (is *ImageSource) Size() int64 {
	return is.size
}

func (is *ImageSource) Cleanup() error {
	if len(is.containerID) > 0 {
		// Remove the temporary container
		if _, err := util.ExecuteCommand("docker", "rm", is.containerID); err != nil {
			return errors.Wrapf(err, "failed to remove container %s:", is.containerID)
		}
	}
	return nil
}

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
	if err := imageFile.Truncate(size + 100*1024*1024); err != nil {
		return errors.Wrapf(err, "failed to allocate space for image %s", md.ID)
	}

	// Use mkfs.ext4 to create the new image with an inode size of 256
	// (gexto doesn't support anything but 128, but as long as we're not using that it's fine)
	if _, err := util.ExecuteCommand("mkfs.ext4", "-I", "256", "-E", "lazy_itable_init=0,lazy_journal_init=0", p); err != nil {
		return errors.Wrapf(err, "failed to format image %s", md.ID)
	}

	return nil
}

// AddFiles copies the contents of the tar file into the ext4 filesystem
func (md *ImageMetadata) AddFiles(src *ImageSource) error {
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
	reader, err := src.GetReader()
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

type KernelNotFoundError struct {
	error
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
		return "", &KernelNotFoundError{err}
	}

	if util.FileExists(kernelSrc) {
		if err := util.CopyFile(kernelSrc, kernelDest); err != nil {
			return "", fmt.Errorf("failed to copy kernel file from %q to %q: %v", kernelSrc, kernelDest, err)
		}
	} else {
		return "", &KernelNotFoundError{fmt.Errorf("no kernel found in image %q", md.ID)}
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
