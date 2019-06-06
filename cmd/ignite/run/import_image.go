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

type ImportImageOptions struct {
	Source     string
	Name       string
	KernelName string
	ImageNames []*metadata.Name
}

func ImportImage(ao *ImportImageOptions) error {
	if !util.FileExists(ao.Source) {
		return fmt.Errorf("not an image file: %s", ao.Source)
	}

	// Create a new ID and directory for the image
	idHandler, err := util.NewID(constants.IMAGE_DIR)
	if err != nil {
		return err
	}
	defer idHandler.Remove()

	// Verify the name
	name, err := metadata.NewName(ao.Name, &ao.ImageNames)
	if err != nil {
		return err
	}

	md := imgmd.NewImageMetadata(idHandler.ID, name)

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

	idHandler.Success()
	return nil
}
