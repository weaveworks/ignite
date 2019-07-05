package operations

import (
	"fmt"
	"log"
	"os"
	"path"

	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
	"github.com/weaveworks/ignite/pkg/source"
	"github.com/weaveworks/ignite/pkg/util"
)

func ImportImage(srcString string, allImages *[]metadata.Metadata) (*imgmd.Image, error) {
	// Parse the source
	dockerSource := source.NewDockerSource()
	src, err := dockerSource.Parse(srcString)
	if err != nil {
		return nil, err
	}

	image := &api.Image{
		Spec: api.ImageSpec{
			Source: *src,
		},
	}

	// Verify the name
	name, err := metadata.NewNameWithLatest(srcString, allImages)
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
	log.Printf("Imported a %s filesystem from OCI image %q", image.Spec.Source.Size.HR(), srcString)
	return runImage, nil
}

// ImportKernelFromImage imports a kernel from an image
// It is assumed that a kernel with the same name does not already exist
// The kernel name will automatically be the same as the image's
// This func returns nil, nil if there is no kernel in the specified image
func ImportKernelFromImage(runImage *imgmd.Image) (*kernmd.Kernel, error) {
	// Import a new kernel from the image if specified
	tmpKernelDir, err := runImage.ExportKernel()
	if err != nil {
		// Tolerate the kernel to not be found
		if _, ok := err.(*imgmd.KernelNotFoundError); ok {
			return nil, nil
		}
		return nil, err
	}

	kernelTmpFile := path.Join(tmpKernelDir, constants.KERNEL_FILE)
	// the kernel name matches the image
	kernelName := runImage.GetName()

	if !util.FileExists(kernelTmpFile) {
		return nil, fmt.Errorf("did not find kernel image: %s", kernelTmpFile)
	}

	// TODO: Kernel importing from docker when moving to pool/snapshotter
	kernel := &api.Kernel{
		Spec: api.KernelSpec{
			Version: "unknown",
			Source: api.ImageSource{
				Type: "file",
				ID:   "-",
				Name: "-",
			},
		},
	}

	// Create new kernel metadata
	runKernel, err := kernmd.NewKernel("", &kernelName, kernel)
	if err != nil {
		return nil, err
	}

	// Save the metadata
	if err := runKernel.Save(); err != nil {
		return nil, err
	}

	// Perform the copy
	filePath := path.Join(runKernel.ObjectPath(), constants.KERNEL_FILE)
	if err := util.CopyFile(kernelTmpFile, filePath); err != nil {
		return nil, fmt.Errorf("failed to copy kernel file %q to kernel %q: %v", kernelTmpFile, runKernel.GetUID(), err)
	}

	// remove the temporary directory
	if err := os.RemoveAll(tmpKernelDir); err != nil {
		return nil, err
	}

	//log.Printf("A kernel was imported from the image with name %q and ID %q", name.String(), kernelID)
	return runKernel, nil
}
