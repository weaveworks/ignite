package run

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/metadata/imgmd"
	"github.com/luxas/ignite/pkg/util"
	"os"
	"path"
)

type ImportImageOptions struct {
	Source     string
	Name       string
	KernelName string
}

func ImportImage(ao *ImportImageOptions) error {
	if !util.FileExists(ao.Source) {
		return fmt.Errorf("not an image file: %s", ao.Source)
	}

	// Create a new ID for the image
	imageID, err := util.NewID(constants.IMAGE_DIR)
	if err != nil {
		return err
	}

	md := imgmd.NewImageMetadata(imageID, ao.Name)

	// Save the metadata
	if err := md.Save(); err != nil {
		return err
	}

	// Perform the copy
	if err := md.ImportImage(ao.Source); err != nil {
		return err
	}

	// Import a new kernel from the image if specified
	if ao.KernelName != "" {
		dir, err := md.ExportKernel()
		if err != nil {
			return err
		}

		if dir != "" {
			if err := ImportKernel(&ImportKernelOptions{
				Source: path.Join(dir, constants.KERNEL_FILE),
				Name:   ao.KernelName,
			}); err != nil {
				return err
			}

			if err := os.RemoveAll(dir); err != nil {
				return err
			}
		}
	}

	fmt.Println(md.ID)

	return nil
}
