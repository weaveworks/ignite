package lookup

import (
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/filter"
)

func ImageUIDForVM(vm *api.VM, c *client.Client) (meta.UID, error) {
	image, err := c.Images().Find(filter.NewNameFilter(vm.Spec.Image.OCIRef.String()))
	if err != nil {
		return "", err
	}

	return image.GetUID(), nil
}

func KernelUIDForVM(vm *api.VM, c *client.Client) (meta.UID, error) {
	kernel, err := c.Kernels().Find(filter.NewNameFilter(vm.Spec.Kernel.OCIRef.String()))
	if err != nil {
		return "", err
	}

	return kernel.GetUID(), nil
}
