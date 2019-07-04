package vmmd

import (
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/metadata"
)

func LoadVM(uid meta.UID) (metadata.Metadata, error) {
	md, err := NewVM(uid, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := md.Load(); err != nil {
		return nil, err
	}

	return md, nil
}

func LoadAllVM() ([]metadata.Metadata, error) {
	return metadata.LoadAllMetadata((&VM{}).TypePath(), LoadVM)
}

func ToVM(md metadata.Metadata) *VM {
	return md.(*VM) // This type assert is internal, we don't need to validate it
}

func ToVMAll(any []metadata.Metadata) []*VM {
	var mds []*VM

	for _, md := range any {
		mds = append(mds, ToVM(md))
	}

	return mds
}
