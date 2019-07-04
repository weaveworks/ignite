package imgmd

import (
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/metadata"
)

func LoadImage(uid meta.UID) (metadata.Metadata, error) {
	md, err := NewImage(uid, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := md.Load(); err != nil {
		return nil, err
	}

	return md, nil
}

func LoadAllImage() ([]metadata.Metadata, error) {
	return metadata.LoadAllMetadata((&Image{}).TypePath(), LoadImage)
}

func ToImage(md metadata.Metadata) *Image {
	return md.(*Image) // This type assert is internal, we don't need to validate it
}

func ToImageAll(any []metadata.Metadata) []*Image {
	var mds []*Image

	for _, md := range any {
		mds = append(mds, ToImage(md))
	}

	return mds
}
