package operations

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/dmlegacy"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/source"
	"github.com/weaveworks/ignite/pkg/util"
	"github.com/weaveworks/libgitops/pkg/filter"
	"github.com/weaveworks/libgitops/pkg/storage/filterer"
)

// FindOrImportImage returns an image based on the source string.
// If the image already exists, it is returned. If the image doesn't
// exist, it is imported
func FindOrImportImage(c *client.Client, ociRef meta.OCIImageRef) (*api.Image, error) {
	log.Debugf("Ensuring image %s exists, or importing it...", ociRef)
	image, err := c.Images().Find(filter.NewIDNameFilter(ociRef.String()))
	if err == nil {
		// Return the image found
		log.Debugf("Found image with UID %s", image.GetUID())
		return image, nil
	}

	switch err.(type) {
	case *filterer.NonexistentError:
		return importImage(c, ociRef)
	default:
		return nil, err
	}
}

// importKernel imports an image from an OCI image
func importImage(c *client.Client, ociRef meta.OCIImageRef) (*api.Image, error) {
	log.Debugf("Importing image with ociRef %q", ociRef)
	// Parse the source
	dockerSource := source.NewDockerSource()
	src, err := dockerSource.Parse(ociRef)
	if err != nil {
		return nil, err
	}

	image := c.Images().New()
	// Set the image name
	image.Name = ociRef.String()
	// Set the image's ociRef
	image.Spec.OCI = ociRef
	// Set the image's ociSource
	image.Status.OCISource = *src

	// Generate UID automatically
	if err := metadata.SetNameAndUID(image, c); err != nil {
		return nil, err
	}

	log.Infoln("Starting image import...")

	// Truncate a file for the filesystem, format it with ext4, and copy in the files from the source
	if err := dmlegacy.CreateImageFilesystem(image, dockerSource); err != nil {
		return nil, err
	}

	if err := c.Images().Set(image); err != nil {
		return nil, err
	}

	log.Infof("Imported OCI image %q (%s) to base image with UID %q", ociRef, image.Status.OCISource.Size, image.GetUID())
	return image, nil
}

// FindOrImportKernel returns an kernel based on the source string.
// If the image already exists, it is returned. If the image doesn't
// exist, it is imported
func FindOrImportKernel(c *client.Client, ociRef meta.OCIImageRef) (*api.Kernel, error) {
	log.Debugf("Ensuring kernel %s exists, or importing it...", ociRef)
	kernel, err := c.Kernels().Find(filter.NewIDNameFilter(ociRef.String()))
	if err == nil {
		// Return the kernel found
		log.Debugf("Found kernel with UID %s", kernel.GetUID())
		return kernel, nil
	}

	switch err.(type) {
	case *filterer.NonexistentError:
		return importKernel(c, ociRef)
	default:
		return nil, err
	}
}

// importKernel imports a kernel from an OCI image
func importKernel(c *client.Client, ociRef meta.OCIImageRef) (*api.Kernel, error) {
	log.Debugf("Importing kernel with ociRef %q", ociRef)
	// Parse the source
	dockerSource := source.NewDockerSource()
	src, err := dockerSource.Parse(ociRef)
	if err != nil {
		return nil, err
	}

	kernel := c.Kernels().New()
	// Set the kernel name
	kernel.Name = ociRef.String()
	// Set the kernel's ociRef
	kernel.Spec.OCI = ociRef
	// Set the kernel's ociSource
	kernel.Status.OCISource = *src

	// Generate UID automatically
	if err := metadata.SetNameAndUID(kernel, c); err != nil {
		return nil, err
	}

	// Cache the kernel contents in the kernel tar file
	kernelTarFile := path.Join(kernel.ObjectPath(), constants.KERNEL_TAR)

	// vmlinuxFile describes the uncompressed kernel file at /var/lib/firecracker/kernel/<id>/vmlinux
	vmlinuxFile := path.Join(kernel.ObjectPath(), constants.KERNEL_FILE)

	// Create both the kernel tar file and the vmlinux file it either doesn't exist
	if !util.FileExists(kernelTarFile) || !util.FileExists(vmlinuxFile) {
		// Create a temporary directory for extracting
		// the necessary files from the OCI image
		tempDir, err := ioutil.TempDir("", "")
		if err != nil {
			return nil, err
		}

		// Extract only the /boot and /lib directories of the tar stream into the tempDir
		err = source.TarExtract(dockerSource, tempDir, "boot", "lib/modules")
		if err != nil {
			return nil, err
		}

		// Locate the kernel file in the temporary directory
		kernelTmpFile, err := findKernel(tempDir)
		if err != nil {
			return nil, err
		}

		// Copy the vmlinux file
		if err := util.CopyFile(kernelTmpFile, vmlinuxFile); err != nil {
			return nil, fmt.Errorf("failed to copy kernel file %q to kernel %q: %v", kernelTmpFile, kernel.GetUID(), err)
		}

		// Pack the kernel tar with unnecessary data removed
		if _, err := util.ExecuteCommand("tar", "-cf", kernelTarFile, "-C", tempDir, "."); err != nil {
			return nil, err
		}

		// Cleanup
		if err := os.RemoveAll(tempDir); err != nil {
			return nil, err
		}
	}

	// Populate the kernel version field if possible
	if len(kernel.Status.Version) == 0 {
		cmd := fmt.Sprintf("strings %s | grep 'Linux version' | awk '{print $3}'", vmlinuxFile)
		// Use the pipefail option to return an error if any of the pipeline commands is not available
		out, err := util.ExecuteCommand("/bin/bash", "-o", "pipefail", "-c", cmd)
		if err != nil {
			kernel.Status.Version = "<unknown>"
		} else {
			kernel.Status.Version = out
		}
	}

	if err := c.Kernels().Set(kernel); err != nil {
		return nil, err
	}

	log.Infof("Imported OCI image %q (%s) to kernel image with UID %q", ociRef, kernel.Status.OCISource.Size, kernel.GetUID())
	return kernel, nil
}

func findKernel(tmpDir string) (string, error) {
	// find the path to the kernel, resolve symlinks if necessary
	bootDir := path.Join(tmpDir, "boot")
	kernel := path.Join(bootDir, constants.KERNEL_FILE)

	fi, err := os.Lstat(kernel)
	if err != nil {
		return "", err
	}

	if fi.Mode()&os.ModeSymlink == 0 {
		// The target file is a real file, not a symlink. Return it
		return kernel, nil
	}

	// The target is a symlink
	kernel, err = os.Readlink(kernel)
	if err != nil {
		return "", err
	}

	// Cleanup the path for absolute and relative symlinks
	if path.IsAbs(kernel) {
		// return the path relative to the tempdir (root)
		// NOTE: This will fail if the symlink starts with any directory other than
		// "/boot", as we don't extract more
		return path.Join(tmpDir, kernel), nil
	}

	// Return the path relative to the boot directory
	return path.Join(bootDir, kernel), nil
}
