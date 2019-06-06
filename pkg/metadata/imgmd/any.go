package imgmd

import (
	"github.com/weaveworks/ignite/pkg/metadata"
)

// Verify that VMMetadata implements AnyMetadata
var _ metadata.AnyMetadata = &ImageMetadata{}

func (md *ImageMetadata) GetMD() *metadata.Metadata {
	return md.Metadata
}

func LoadImageMetadata(id string) (metadata.AnyMetadata, error) {
	md := NewImageMetadata(id, nil)
	err := md.Load()
	return md, err
}

func LoadAllImageMetadata() ([]metadata.AnyMetadata, error) {
	return metadata.LoadAllMetadata(metadata.Image.Path(), LoadImageMetadata)
}

func ToImageMetadata(md metadata.AnyMetadata) *ImageMetadata {
	return md.(*ImageMetadata) // This type assert is internal, we don't need to validate it
}

func ToImageMetadataAll(any []metadata.AnyMetadata) []*ImageMetadata {
	var mds []*ImageMetadata

	for _, md := range any {
		mds = append(mds, ToImageMetadata(md))
	}

	return mds
}
