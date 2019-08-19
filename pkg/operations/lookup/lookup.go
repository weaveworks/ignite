package lookup

import (
	"github.com/weaveworks/gitops-toolkit/pkg/filter"
	"github.com/weaveworks/gitops-toolkit/pkg/runtime"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/client"
)

func ImageUIDForVM(vm *api.VM, c *client.Client) (runtime.UID, error) {
	image, err := c.Images().Find(filter.NewNameFilter(vm.Spec.Image.OCI.String()))
	if err != nil {
		return "", err
	}

	return image.GetUID(), nil
}

func KernelUIDForVM(vm *api.VM, c *client.Client) (runtime.UID, error) {
	kernel, err := c.Kernels().Find(filter.NewNameFilter(vm.Spec.Kernel.OCI.String()))
	if err != nil {
		return "", err
	}

	return kernel.GetUID(), nil
}
