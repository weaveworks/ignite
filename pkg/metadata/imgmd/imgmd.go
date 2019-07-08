package imgmd

import (
	"path"

	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata"
)

type Image struct {
	*api.Image
}

var _ metadata.Metadata = &Image{}

func NewImage(id meta.UID, name *string, object *api.Image) (*Image, error) {
	if object == nil {
		object = &api.Image{}
	}
	// Set defaults, and populate TypeMeta
	// TODO: Make this more standardized; maybe a constructor method somewhere?
	scheme.Scheme.Default(object)

	md := &Image{
		Image: object,
	}

	metadata.InitName(md, name)

	if err := metadata.NewUID(md, id); err != nil {
		return nil, err
	}

	return md, nil
}

// TODO: Remove
func (md *Image) TypePath() string {
	return constants.IMAGE_DIR
}

func (md *Image) ObjectPath() string {
	return path.Join(md.TypePath(), md.GetUID().String())
}

func (md *Image) Load() (err error) {
	md.Image, err = client.Images().Get(md.GetUID())
	return
}

func (md *Image) Save() error {
	return client.Images().Set(md.Image)
}
