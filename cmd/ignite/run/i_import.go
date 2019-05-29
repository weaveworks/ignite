package run

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/metadata/imgmd"
	"github.com/luxas/ignite/pkg/util"
)

type ImportImageOptions struct {
	Source string
	Name   string
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

	fmt.Println(md.ID)

	return nil
}
