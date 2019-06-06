package imgmd

import (
	"github.com/weaveworks/ignite/pkg/metadata"
)

type ImageMetadata struct {
	*metadata.Metadata
}

type ImageObjectData struct {
	// TODO: Placeholder
}

func NewImageMetadata(id string, name *metadata.Name) *ImageMetadata {
	return &ImageMetadata{
		Metadata: metadata.NewMetadata(
			id,
			name,
			metadata.Image,
			&ImageObjectData{}),
	}
}

// The md.ObjectData.(*ImageObjectData) assert won't panic as this method can only receive *ImageMetadata objects
func (md *ImageMetadata) ImageOD() *ImageObjectData {
	return md.ObjectData.(*ImageObjectData)
}
