package run

import (
	"log"
	"os"
	"path"

	"github.com/c2h5oh/datasize"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/util"
)

type BuildOptions struct {
	Source     string
	Name       string
	KernelName string
	image      *imgmd.ImageMetadata
	ImageNames []*metadata.Name
}

func Build(bo *BuildOptions) (string, error) {
	// Create a new ID and directory for the image
	idHandler, err := util.NewID(constants.IMAGE_DIR)
	if err != nil {
		return "", err
	}
	defer idHandler.Remove()

	// Parse the source
	imageSrc, err := imgmd.NewSource(bo.Source)
	if err != nil {
		return "", err
	}

	nameStr := bo.Name
	if len(imageSrc.DockerImage()) > 0 {
		nameStr = imageSrc.DockerImage()
	}

	// Verify the name
	name, err := metadata.NewNameWithLatest(nameStr, &bo.ImageNames)
	if err != nil {
		return "", err
	}

	// Create new image metadata
	bo.image = imgmd.NewImageMetadata(idHandler.ID, name)

	log.Println("Starting image build...")

	// Create new file to host the filesystem and format it
	if err := bo.image.AllocateAndFormat(imageSrc.Size()); err != nil {
		return "", err
	}

	// Add the files to the filesystem
	if err := bo.image.AddFiles(imageSrc); err != nil {
		return "", err
	}

	if err := bo.image.Save(); err != nil {
		return "", err
	}
	hrsize := datasize.ByteSize(imageSrc.Size()).HR()
	log.Printf("Created a %s filesystem of the input", hrsize)

	// Import a new kernel from the image if specified
	tmpKernelDir, err := bo.image.ExportKernel()
	if err == nil {
		_, err := ImportKernel(&ImportKernelOptions{
			Source: path.Join(tmpKernelDir, constants.KERNEL_FILE),
			Name:   name.String(),
		})
		if err != nil {
			return "", err
		}
		if err := os.RemoveAll(tmpKernelDir); err != nil {
			return "", err
		}
		//log.Printf("A kernel was imported from the image with name %q and ID %q", name.String(), kernelID)
	} else {
		// Tolerate the kernel to not be found
		if _, ok := err.(*imgmd.KernelNotFoundError); !ok {
			return "", err
		}
	}
	return idHandler.Success(name.String())
}
