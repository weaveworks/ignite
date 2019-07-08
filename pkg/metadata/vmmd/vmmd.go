package vmmd

import (
	"path"

	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"

	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata"
)

type VM struct {
	*api.VM
}

var _ metadata.Metadata = &VM{}

func NewVM(id meta.UID, name *string, object *api.VM) (*VM, error) {
	if object == nil {
		object = &api.VM{}
	}

	md := &VM{
		VM: object,
	}

	metadata.InitName(md, name)

	if err := metadata.NewUID(md, id); err != nil {
		return nil, err
	}

	return md, nil
}

// TODO: Remove
func (md *VM) TypePath() string {
	return constants.VM_DIR
}

func (md *VM) ObjectPath() string {
	return path.Join(md.TypePath(), md.GetUID().String())
}

func (md *VM) Load() (err error) {
	md.VM, err = client.VMs().Get(md.GetUID())
	return
}

func (md *VM) Save() error {
	return client.VMs().Set(md.VM)
}
