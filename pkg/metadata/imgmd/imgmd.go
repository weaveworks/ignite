package imgmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/metadata"
)

type ImageMetadata struct {
	*metadata.Metadata
}

type ImageObjectData struct {
	// TODO: Placeholder
}

func NewImageMetadata(id, name string) *ImageMetadata {
	return &ImageMetadata{
		Metadata: &metadata.Metadata{
			ID:         id,
			Name:       name,
			Type:       metadata.Image,
			ObjectData: &ImageObjectData{},
		},
	}
}

func ToImageMetadata(f filter.Filterable) (*ImageMetadata, error) {
	md, ok := f.(*ImageMetadata)
	if !ok {
		return nil, fmt.Errorf("failed to assert Filterable %v to ImageMetadata", f)
	}

	return md, nil
}

func ToImageMetadataAll(a []filter.Filterable) ([]*ImageMetadata, error) {
	var mds []*ImageMetadata

	for _, f := range a {
		if md, err := ToImageMetadata(f); err == nil {
			mds = append(mds, md)
		} else {
			return nil, err
		}
	}

	return mds, nil
}

// The md.ObjectData.(*ImageObjectData) assert won't panic as these methods can only receive *KernelMetadata objects
//func (md *ImageMetadata) Placeholder() error {
//	od := md.ObjectData.(*ImageObjectData)
//
//	return nil
//}
