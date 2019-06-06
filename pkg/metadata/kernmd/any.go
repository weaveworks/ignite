package kernmd

import (
	"github.com/weaveworks/ignite/pkg/metadata"
)

// Verify that VMMetadata implements AnyMetadata
var _ metadata.AnyMetadata = &KernelMetadata{}

func (md *KernelMetadata) GetMD() *metadata.Metadata {
	return md.Metadata
}

func LoadKernelMetadata(id string) (metadata.AnyMetadata, error) {
	md := NewKernelMetadata(id, nil)
	err := md.Load()
	return md, err
}

func LoadAllKernelMetadata() ([]metadata.AnyMetadata, error) {
	return metadata.LoadAllMetadata(metadata.Kernel.Path(), LoadKernelMetadata)
}

func ToKernelMetadata(md metadata.AnyMetadata) *KernelMetadata {
	return md.(*KernelMetadata) // This type assert is internal, we don't need to validate it
}

func ToKernelMetadataAll(any []metadata.AnyMetadata) []*KernelMetadata {
	var mds []*KernelMetadata

	for _, md := range any {
		mds = append(mds, ToKernelMetadata(md))
	}

	return mds
}
