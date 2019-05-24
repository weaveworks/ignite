package kernmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/metadata"
)

type KernelMetadata struct {
	*metadata.Metadata
}

type KernelObjectData struct {
	// TODO: Placeholder
}

func NewKernelMetadata(id, name string) *KernelMetadata {
	return &KernelMetadata{
		Metadata: &metadata.Metadata{
			ID:         id,
			Name:       name,
			Type:       metadata.Kernel,
			ObjectData: &KernelObjectData{},
		},
	}
}

func ToKernelMetadata(f filter.Filterable) (*KernelMetadata, error) {
	md, ok := f.(*KernelMetadata)
	if !ok {
		return nil, fmt.Errorf("failed to assert Filterable %v to KernelMetadata", f)
	}

	return md, nil
}

func ToKernelMetadataAll(a []filter.Filterable) ([]*KernelMetadata, error) {
	var mds []*KernelMetadata

	for _, f := range a {
		if md, err := ToKernelMetadata(f); err == nil {
			mds = append(mds, md)
		} else {
			return nil, err
		}
	}

	return mds, nil
}

// The md.ObjectData.(*KernelObjectData) assert won't panic as these methods can only receive *KernelMetadata objects
//func (md *KernelMetadata) Placeholder() error {
//	od := md.ObjectData.(*KernelObjectData)
//
//	return nil
//}
