package vmmd

import (
	"github.com/weaveworks/ignite/pkg/metadata"
)

func LoadVMMetadata(id string) (metadata.Metadata, error) {
	md, err := NewVMMetadata(id, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := md.Load(); err != nil {
		return nil, err
	}

	return md, nil
}

func LoadAllVMMetadata() ([]metadata.Metadata, error) {
	return metadata.LoadAllMetadata((&VMMetadata{}).TypePath(), LoadVMMetadata)
}

func ToVMMetadata(md metadata.Metadata) *VMMetadata {
	return md.(*VMMetadata) // This type assert is internal, we don't need to validate it
}

func ToVMMetadataAll(any []metadata.Metadata) []*VMMetadata {
	var mds []*VMMetadata

	for _, md := range any {
		mds = append(mds, ToVMMetadata(md))
	}

	return mds
}
