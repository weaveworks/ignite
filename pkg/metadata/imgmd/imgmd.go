package imgmd

import (
	"github.com/weaveworks/ignite/pkg/metadata"
)

type ImageMetadata struct {
	*metadata.Metadata
}

type ImageObjectData struct {
	ImageDM *ImageDM
}

func NewImageMetadata(id *metadata.ID, name *metadata.Name) (*ImageMetadata, error) {
	md, err := metadata.NewMetadata(id, name, metadata.Image, NewImageObjectData())
	if err != nil {
		return nil, err
	}

	return &ImageMetadata{Metadata: md}, nil
}

func NewImageObjectData() *ImageObjectData {
	return &ImageObjectData{
		ImageDM: emptyImageDM(),
	}
}

// The md.ObjectData.(*ImageObjectData) assert won't panic as this method can only receive *ImageMetadata objects
func (md *ImageMetadata) ImageOD() *ImageObjectData {
	return md.ObjectData.(*ImageObjectData)
}
