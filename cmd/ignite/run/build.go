package run

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/metadata"
	"github.com/luxas/ignite/pkg/metadata/imgmd"
	"github.com/luxas/ignite/pkg/util"
	"os"
	"path"
)

type BuildOptions struct {
	Source     string
	Name       string
	KernelName string
	image      *imgmd.ImageMetadata
	ImageNames []*metadata.Name
}

func Build(bo *BuildOptions) error {
	// Create new image metadata
	if err := bo.newImage(); err != nil {
		return err
	}

	imageSrc, err := imgmd.NewSource(bo.Source)
	if err != nil {
		return err
	}

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
	if bo.KernelName != "" {
		dir, err := bo.image.ExportKernel()
		if err != nil {
			return err
		}

		if dir != "" {
			if err := ImportKernel(&ImportKernelOptions{
				Source: path.Join(dir, constants.KERNEL_FILE),
				Name:   bo.KernelName,
			}); err != nil {
				return err
			}

			if err := os.RemoveAll(dir); err != nil {
				return err
			}
		}
	}

	//if err := container.ExportToDocker(image); err != nil {
	//	return err
	//}

	// Print the ID of the newly generated image
	fmt.Println(bo.image.ID)

	return nil
}

// newImage creates the new image metadata
func (bo *BuildOptions) newImage() error {
	newID, err := util.NewID(constants.IMAGE_DIR)
	if err != nil {
		return err
	}

	// Verify the name
	name, err := metadata.NewName(bo.Name, &bo.ImageNames)
	if err != nil {
		return err
	}

	bo.image = imgmd.NewImageMetadata(newID, name)

	return nil
}
