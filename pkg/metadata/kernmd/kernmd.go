package kernmd

import (
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/client"
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
	return &Kernel{
		Kernel: obj,
		c:      client.DefaultClient,
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
