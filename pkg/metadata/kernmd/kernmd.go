package kernmd

import (
	"path"

	"github.com/weaveworks/ignite/pkg/providers"

	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata"
)

type Kernel struct {
	*api.Kernel
	c *client.Client
}

var _ metadata.Metadata = &Kernel{}

// WrapKernel wraps an API type in the runtime object
// It does not do any validation or checking like
// NewKernel, hence it should only be used for "safe"
// data coming from storage.
func WrapKernel(obj *api.Kernel) *Kernel {
	// Run the object through defaulting, just to be sure it has all the values
	scheme.Scheme.Default(obj)

	return &Kernel{
		Kernel: obj,
		c:      providers.Client,
	}
}

func NewKernel(obj *api.Kernel, c *client.Client) (*Kernel, error) {
	// Initialize UID, name, defaulting, etc. that is common for all kinds
	if err := metadata.InitObject(obj, c); err != nil {
		return nil, err
	}

	// TODO: Validate the API object here

	// Construct the runtime object
	kernel := &Kernel{
		Kernel: obj,
		c:      c,
	}

	return kernel, nil
}

func (k *Kernel) ObjectPath() string {
	// TODO: Move this into storage
	return path.Join(constants.DATA_DIR, k.GetKind().Lower(), k.GetUID().String())
}

func (k *Kernel) Save() error {
	return k.c.Kernels().Set(k.Kernel)
}
