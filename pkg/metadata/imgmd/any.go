package imgmd

import (
	"github.com/weaveworks/ignite/pkg/metadata"
)

func LoadImageMetadata(id string) (metadata.Metadata, error) {
	md, err := NewImageMetadata(id, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := md.Load(); err != nil {
		return nil, err
	}

	return md, nil
}

func LoadAllImageMetadata() ([]metadata.Metadata, error) {
	return metadata.LoadAllMetadata((&ImageMetadata{}).TypePath(), LoadImageMetadata)
}

func ToImageMetadata(md metadata.Metadata) *ImageMetadata {
	return md.(*ImageMetadata) // This type assert is internal, we don't need to validate it
}

func ToImageMetadataAll(any []metadata.Metadata) []*ImageMetadata {
	var mds []*ImageMetadata

	for _, md := range any {
		mds = append(mds, ToImageMetadata(md))
	}

	return mds
}
