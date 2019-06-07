package run

import (
	"fmt"
	"os"
	"path"

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

func Build(bo *BuildOptions) error {
	// Create a new ID and directory for the image
	idHandler, err := util.NewID(constants.IMAGE_DIR)
	if err != nil {
		return err
	}
	defer idHandler.Remove()

	// Parse the source
	imageSrc, err := imgmd.NewSource(bo.Source)
	if err != nil {
		return err
	}

	nameStr := bo.Name
	if len(imageSrc.DockerImage()) > 0 {
		nameStr = imageSrc.DockerImage()
	}

	// Verify the name
	name, err := metadata.NewNameWithLatest(nameStr, &bo.ImageNames)
	if err != nil {
		return err
	}

	// Create new image metadata
	bo.image = imgmd.NewImageMetadata(idHandler.ID, name)

	// Create new file to host the filesystem and format it
	if err := bo.image.AllocateAndFormat(imageSrc.Size()); err != nil {
		return err
	}

	// Add the files to the filesystem
	if err := bo.image.AddFiles(imageSrc); err != nil {
		return err
	}

	if err := bo.image.Save(); err != nil {
		return err
	}

	// Import a new kernel from the image if specified
	dir, err := bo.image.ExportKernel()
	if err == nil {
		if dir != "" {
			if err := ImportKernel(&ImportKernelOptions{
				Source: path.Join(dir, constants.KERNEL_FILE),
				Name:   name.String(),
			}); err != nil {
				return err
			}

			if err := os.RemoveAll(dir); err != nil {
				return err
			}
		}
	} else {
		// Tolerate the kernel to not be found
		if _, ok := err.(*imgmd.KernelNotFoundError); !ok {
			return err
		}
	}

	// Print the ID of the newly generated image
	fmt.Println(bo.image.ID)

	idHandler.Success()
	return nil
}
