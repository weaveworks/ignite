package operations

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/filter"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
	"github.com/weaveworks/ignite/pkg/source"
	"github.com/weaveworks/ignite/pkg/storage/filterer"
	"github.com/weaveworks/ignite/pkg/util"
)

// FindOrImportImage returns an image based on the source string.
// If the image already exists, it is returned. If the image doesn't
// exist, it is imported
func FindOrImportImage(c *client.Client, srcString string) (*imgmd.Image, error) {
	image, err := c.Images().Find(filter.NewIDNameFilter(srcString))
	if err == nil {
		return &imgmd.Image{image}, nil
	}

	switch err.(type) {
	case *filterer.NonexistentError:
		return importImage(srcString)
	default:
		return nil, err
	}
}

// importKernel imports an image from an OCI image
func importImage(srcString string) (*imgmd.Image, error) {
	// Parse the source
	dockerSource := source.NewDockerSource()
	src, err := dockerSource.Parse(srcString)
	if err != nil {
		return nil, err
	}

	image := &api.Image{
		Spec: api.ImageSpec{
			OCIClaim: api.OCIImageClaim{
				Ref: srcString,
			},
		},
		Status: api.ImageStatus{
			OCISource: *src,
		},
	}
	// Set defaults, and populate TypeMeta
	// TODO: Make this more standardized; maybe a constructor method somewhere?
	scheme.Scheme.Default(image)

	// Verify the name
	name, err := metadata.NewNameWithLatest(srcString, meta.KindImage)
	if err != nil {
		return nil, err
	}

	// Create new image metadata
	runImage, err := imgmd.NewImage("", &name, image)
	if err != nil {
		return nil, err
	}

	log.Println("Starting image import...")

	// Create new file to host the filesystem and format it
	if err := runImage.AllocateAndFormat(); err != nil {
		return nil, err
	}

	// Add the files to the filesystem
	if err := runImage.AddFiles(dockerSource); err != nil {
		return nil, err
	}

	if err := runImage.Save(); err != nil {
		return nil, err
	}
	log.Printf("Imported a %s filesystem from OCI image %q", image.Status.OCISource.Size.HR(), srcString)
	return runImage, nil
}

// FindOrImportKernel returns an kernel based on the source string.
// If the image already exists, it is returned. If the image doesn't
// exist, it is imported
func FindOrImportKernel(c *client.Client, srcString string) (*kernmd.Kernel, error) {
	image, err := c.Kernels().Find(filter.NewIDNameFilter(srcString))
	if err == nil {
		return &kernmd.Kernel{image}, nil
	}

	switch err.(type) {
	case *filterer.NonexistentError:
		return importKernel(srcString)
	default:
		return nil, err
	}
}

// importKernel imports a kernel from an OCI image
func importKernel(srcString string) (*kernmd.Kernel, error) {
	// Parse the source
	dockerSource := source.NewDockerSource()
	src, err := dockerSource.Parse(srcString)
	if err != nil {
		return nil, err
	}

	kernel := &api.Kernel{
		Spec: api.KernelSpec{
			OCIClaim: api.OCIImageClaim{
				Ref: srcString,
			},
		},
		Status: api.KernelStatus{
			OCISource: *src,
		},
	}

	// Set defaults, and populate TypeMeta
	// TODO: Make this more standardized; maybe a constructor method somewhere?
	scheme.Scheme.Default(kernel)

	name, err := metadata.NewNameWithLatest(srcString, meta.KindKernel)
	if err != nil {
		return nil, err
	}

	// Create new kernel metadata
	runKernel, err := kernmd.NewKernel("", &name, kernel)
	if err != nil {
		return nil, err
	}

	// Cache the kernel contents in the kernel tar file
	kernelTarFile := path.Join(runKernel.ObjectPath(), constants.KERNEL_TAR)

	if !util.FileExists(kernelTarFile) {
		f, err := os.Create(kernelTarFile)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		reader, err := dockerSource.Reader()
		if err != nil {
			return nil, err
		}
		defer reader.Close()

		// Copy over the contents from the OCI image into the tar file
		if _, err := io.Copy(f, reader); err != nil {
			return nil, err
		}

		// Remove the temporary container
		if err := dockerSource.Cleanup(); err != nil {
			return nil, err
		}
	}

	// vmlinuxFile describes the uncompressed kernel file at /var/lib/firecracker/kernel/$id/vmlinux
	vmlinuxFile := path.Join(runKernel.ObjectPath(), constants.KERNEL_FILE)
	// Create it if it doesn't exist
	if !util.FileExists(vmlinuxFile) {
		// Create a temporary directory for extracting the kernel file
		tempDir, err := ioutil.TempDir("", "")
		if err != nil {
			return nil, err
		}
		// Extract only the boot directory from the tar file cache to the temp dir
		if _, err := util.ExecuteCommand("tar", "-xf", kernelTarFile, "-C", tempDir, "boot"); err != nil {
			return nil, err
		}

		// Locate the kernel file in the temporary directory
		kernelTmpFile, err := findKernel(tempDir)
		if err != nil {
			return nil, err
		}

		// Copy the vmlinux file
		if err := util.CopyFile(kernelTmpFile, vmlinuxFile); err != nil {
			return nil, fmt.Errorf("failed to copy kernel file %q to kernel %q: %v", kernelTmpFile, runKernel.GetUID(), err)
		}

		// Cleanup
		if err := os.RemoveAll(tempDir); err != nil {
			return nil, err
		}
	}

	// Populate the kernel version field if possible
	if len(runKernel.Status.Version) == 0 {
		cmd := fmt.Sprintf(`strings %s | grep 'Linux version' | awk '{print $3}'`, vmlinuxFile)
		out, err := util.ExecuteCommand("/bin/bash", "-c", cmd)
		if err != nil {
			runKernel.Status.Version = "<unknown>"
		} else {
			runKernel.Status.Version = string(out)
		}
	}

	// Save the metadata
	if err := runKernel.Save(); err != nil {
		return nil, err
	}

	log.Printf("A kernel was imported from the image with name %q and ID %q", runKernel.GetName(), runKernel.GetUID())
	return runKernel, nil
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
