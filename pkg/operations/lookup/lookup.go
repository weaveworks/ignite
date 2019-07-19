package lookup

import (
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/filter"
)

func ImageUIDForVM(vm *api.VM, c *client.Client) (meta.UID, error) {
	image, err := c.Images().Find(filter.NewNameFilter(vm.Spec.Image.OCIClaim.Ref.String()))
	if err != nil {
		return meta.UID(""), err
	}
	return image.GetUID(), nil
}

func KernelUIDForVM(vm *api.VM, c *client.Client) (meta.UID, error) {
	kernel, err := c.Images().Find(filter.NewNameFilter(vm.Spec.Kernel.OCIClaim.Ref.String()))
	if err != nil {
		return meta.UID(""), err
	}
	return kernel.GetUID(), nil
}
