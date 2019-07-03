package vmmd

import (
	"path"

	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata"
)

type VMMetadata struct {
	*v1alpha1.VM
}

var _ metadata.Metadata = &VMMetadata{}

func NewVMMetadata(id string, name *string, object *v1alpha1.VM) (*VMMetadata, error) {
	if object == nil {
		object = &v1alpha1.VM{}
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

func (md *VMMetadata) Type() v1alpha1.PoolDeviceType {
	return v1alpha1.PoolDeviceTypeVM
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
