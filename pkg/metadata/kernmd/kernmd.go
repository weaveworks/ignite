package kernmd

import (
	"path"

	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata"
)

type Kernel struct {
	*api.Kernel
}

var _ metadata.Metadata = &Kernel{}

func NewKernel(id string, name *string, object *api.Kernel) (*Kernel, error) {
	if object == nil {
		object = &api.Kernel{}
	}

	md := &Kernel{
		Kernel: object,
	}

	metadata.InitName(md, name)

	if err := metadata.NewID(md, id); err != nil {
		return nil, err
	}

	return md, nil
}

func (md *Kernel) Type() api.PoolDeviceType {
	return api.PoolDeviceTypeKernel
}

func (md *Kernel) TypePath() string {
	return constants.KERNEL_DIR
}

func (md *Kernel) ObjectPath() string {
	return path.Join(md.TypePath(), md.GetUID())
}

func (md *Kernel) Load() (err error) {
	md.Kernel, err = client.Kernels().Get(md.GetUID())
	return
}

func (md *Kernel) Save() error {
	return client.Kernels().Set(md.Kernel)
}
