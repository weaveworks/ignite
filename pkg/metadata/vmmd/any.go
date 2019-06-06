package vmmd

import (
	"github.com/weaveworks/ignite/pkg/metadata"
)

// Verify that VMMetadata implements AnyMetadata
var _ metadata.AnyMetadata = &VMMetadata{}

func (md *VMMetadata) GetMD() *metadata.Metadata {
	return md.Metadata
}

func LoadVMMetadata(id string) (metadata.AnyMetadata, error) {
	md := NewVMMetadata(id, nil, &VMObjectData{})
	err := md.Load()
	return md, err
}

func LoadAllVMMetadata() ([]metadata.AnyMetadata, error) {
	return metadata.LoadAllMetadata(metadata.VM.Path(), LoadVMMetadata)
}

func ToVMMetadata(md metadata.AnyMetadata) *VMMetadata {
	return md.(*VMMetadata) // This type assert is internal, we don't need to validate it
}

func ToVMMetadataAll(any []metadata.AnyMetadata) []*VMMetadata {
	var mds []*VMMetadata

	for _, md := range any {
		mds = append(mds, ToVMMetadata(md))
	}

	return mds
}
