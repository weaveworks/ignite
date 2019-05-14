package build

import (
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/util"
	"github.com/pkg/errors"
	"os"
	"path"
)

type Image struct {
	id   string
	path string
}

func NewImage(id string) *Image {
	return &Image{
		id:   id,
		path: path.Join(constants.IMAGE_DIR, id, constants.IMAGE_FS),
	}
}

func (i Image) AllocateAndFormat() error {
	imageFile, err := os.Create(i.path)
	if err != nil {
		return errors.Wrapf(err, "failed to create image file for %s", i.id)
	}

	// TODO: Dynamic size, for now hardcoded 4 GiB
	if err := imageFile.Truncate(4294967296); err != nil {
		return errors.Wrapf(err, "failed to allocate space for image %s", i.id)
	}

	if _, err := util.ExecuteCommand("mkfs.ext4", i.path); err != nil {
		return errors.Wrapf(err, "failed to format image %s", i.id)
	}

	return nil
}
