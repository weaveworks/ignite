package snapshotter

import (
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/dm"
)

type VM struct {
	v1alpha1.VM
	layer
}

func NewVM(vm v1alpha1.VM, device *dm.Device) *VM {
	return &VM{
		VM:    vm,
		layer: newLayer(device),
	}
}
