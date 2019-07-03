package kernmd

import (
	"path"

	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata"
)

type KernelMetadata struct {
	*api.Kernel
}

var _ metadata.Metadata = &KernelMetadata{}

func NewKernelMetadata(id string, name *string, object *api.Kernel) (*KernelMetadata, error) {
	if object == nil {
		object = &api.Kernel{}
	}

	md := &KernelMetadata{
		Kernel: object,
	}

	metadata.InitName(md, name)

	if err := metadata.NewID(md, id); err != nil {
		return nil, err
	}

	return md, nil
}

func (md *KernelMetadata) Type() api.PoolDeviceType {
	return api.PoolDeviceTypeKernel
}

func (md *KernelMetadata) TypePath() string {
	return constants.KERNEL_DIR
}

func (md *KernelMetadata) ObjectPath() string {
	return path.Join(md.TypePath(), md.GetUID())
}

func (md *KernelMetadata) Load() (err error) {
	md.Kernel, err = client.Kernels().Get(md.GetUID())
	return
}

func (md *KernelMetadata) Save() error {
	return client.Kernels().Set(md.Kernel)
}
