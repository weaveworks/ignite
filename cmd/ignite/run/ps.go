package run

import (
	"fmt"
	"github.com/luxas/ignite/pkg/metadata/imgmd"
	"github.com/luxas/ignite/pkg/metadata/kernmd"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"github.com/luxas/ignite/pkg/util"
)

type PsOptions struct {
	VMs []*vmmd.VMMetadata
	All bool
}

func Ps(po *PsOptions) error {
	o := util.NewOutput()
	defer o.Flush()

	o.Write("VM ID", "IMAGE", "KERNEL", "CREATED", "SIZE", "CPUS", "MEMORY", "STATE", "NAME")
	for _, vm := range po.VMs {
		od := vm.ObjectData.(*vmmd.VMObjectData)
		size, err := vm.Size()
		if err != nil {
			return fmt.Errorf("failed to get size for %s %q: %v", vm.Type, vm.ID, err)
		}

		image, err := imgmd.LoadImageMetadata(od.ImageID)
		if err != nil {
			return fmt.Errorf("failed to load image metadata for %s %q: %v", vm.Type, vm.ID, err)
		}

		kernel, err := kernmd.LoadKernelMetadata(od.KernelID)
		if err != nil {
			return fmt.Errorf("failed to load kernel metadata for %s %q: %v", vm.Type, vm.ID, err)
		}

		// TODO: Clean up this print
		o.Write(vm.ID, image.Name, kernel.Name, vm.Created, util.ByteCountDecimal(size), od.VCPUs, util.ByteCountDecimal(od.Memory*1000000), od.State, vm.Name)
	}

	return nil
}
