package imgmd

import (
	"path"

	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata"
)

type Image struct {
	*api.Image
}

var _ metadata.Metadata = &Image{}

func NewImage(id string, name *string, object *api.Image) (*Image, error) {
	if object == nil {
		object = &api.Image{}
	}

	md := &Image{
		Image: object,
	}

	metadata.InitName(md, name)

	if err := metadata.NewID(md, id); err != nil {
		return nil, err
	}

	return md, nil
}

func (md *Image) Type() api.PoolDeviceType {
	return api.PoolDeviceTypeImage
}

func (md *Image) TypePath() string {
	return constants.IMAGE_DIR
}

func (md *Image) ObjectPath() string {
	return path.Join(md.TypePath(), md.GetUID())
}

func (md *Image) Load() (err error) {
	md.Image, err = client.Images().Get(md.GetUID())
	return
}

func (md *Image) Save() error {
	return client.Images().Set(md.Image)
}
