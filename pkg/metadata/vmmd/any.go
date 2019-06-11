package vmmd

import (
	"github.com/weaveworks/ignite/pkg/metadata"
)

// Verify that VMMetadata implements AnyMetadata
var _ metadata.AnyMetadata = &VMMetadata{}

func (md *VMMetadata) GetMD() *metadata.Metadata {
	return md.Metadata
}

func LoadVMMetadata(id *metadata.ID) (metadata.AnyMetadata, error) {
	md, err := NewVMMetadata(id, nil, &VMObjectData{})
	if err != nil {
		return nil, err
	}

	if err := md.Load(); err != nil {
		return nil, err
	}

	return md, nil
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
