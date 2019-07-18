package imgmd

import (
	"path"

	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata"
)

type Image struct {
	*api.Image
	c *client.Client
}

var _ metadata.Metadata = &Image{}

// WrapImage wraps an API type in the runtime object
// It does not do any validation or checking like
// NewImage, hence it should only be used for "safe"
// data coming from storage.
func WrapImage(obj *api.Image) *Image {
	// Run the object through defaulting, just to be sure it has all the values
	scheme.Serializer.DefaultInternal(obj)

	return &Image{
		Image: obj,
		c:     client.DefaultClient,
	}
}

func NewImage(obj *api.Image, c *client.Client) (*Image, error) {
	// Initialize UID, name, defaulting, etc. that is common for all kinds
	if err := metadata.InitObject(obj, c); err != nil {
		return nil, err
	}

	// TODO: Validate the API object here

	// Construct the runtime object
	md := &Image{
		Image: obj,
		c:     c,
	}
	return md, nil
}

func (img *Image) ObjectPath() string {
	// TODO: Move this into storage
	return path.Join(constants.DATA_DIR, img.GetKind().Lower(), img.GetUID().String())
}

func (img *Image) Save() error {
	return img.c.Images().Set(img.Image)
}
