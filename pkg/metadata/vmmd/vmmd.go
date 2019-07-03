package vmmd

import (
	"path"

	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata"
)

type VMMetadata struct {
	*api.VM
}

var _ metadata.Metadata = &VMMetadata{}

func NewVMMetadata(id string, name *string, object *api.VM) (*VMMetadata, error) {
	if object == nil {
		object = &api.VM{}
	}

	md := &VMMetadata{
		VM: object,
	}

	metadata.InitName(md, name)

	if err := metadata.NewID(md, id); err != nil {
		return nil, err
	}

	return md, nil
}

func (md *VMMetadata) Type() api.PoolDeviceType {
	return api.PoolDeviceTypeVM
}

func (md *VMMetadata) TypePath() string {
	return constants.VM_DIR
}

func (md *VMMetadata) ObjectPath() string {
	return path.Join(md.TypePath(), md.GetUID())
}

func (md *VMMetadata) Load() (err error) {
	md.VM, err = client.VMs().Get(md.GetUID())
	return
}

func (md *VMMetadata) Save() error {
	return client.VMs().Set(md.VM)
}
