package kernmd

import (
	"github.com/weaveworks/ignite/pkg/metadata"
)

func LoadKernelMetadata(id string) (metadata.Metadata, error) {
	md, err := NewKernelMetadata(id, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := md.Load(); err != nil {
		return nil, err
	}

	return md, nil
}

func LoadAllKernelMetadata() ([]metadata.Metadata, error) {
	return metadata.LoadAllMetadata((&KernelMetadata{}).TypePath(), LoadKernelMetadata)
}

func ToKernelMetadata(md metadata.Metadata) *KernelMetadata {
	return md.(*KernelMetadata) // This type assert is internal, we don't need to validate it
}

func ToKernelMetadataAll(any []metadata.Metadata) []*KernelMetadata {
	var mds []*KernelMetadata

	for _, md := range any {
		mds = append(mds, ToKernelMetadata(md))
	}

	return mds
}
