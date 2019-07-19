package vmmd

import (
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/metadata"
)

type VM struct {
	*api.VM
}

var _ metadata.Metadata = &VM{}

// WrapVM wraps an API type in the runtime object
// It does not do any validation or checking like
// NewVM, hence it should only be used for "safe"
// data coming from storage.
func WrapVM(obj *api.VM) *VM {
	vm := &VM{
		VM: obj,
	}

	return vm
}

func NewVM(obj *api.VM, c *client.Client) (*VM, error) {
	// Initialize UID, name, defaulting, etc. that is common for all kinds
	if err := metadata.SetNameAndUID(obj, c); err != nil {
		return nil, err
	}

	// Construct the runtime object
	vm := &VM{
		VM: obj,
	}

	return vm, nil
}
