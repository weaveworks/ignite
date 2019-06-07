package run

import (
	"fmt"
	"strings"

	"github.com/c2h5oh/datasize"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/util"
)

type PsOptions struct {
	VMs []*vmmd.VMMetadata
	All bool
}

func Ps(po *PsOptions) error {
	o := util.NewOutput()
	defer o.Flush()

	o.Write("VM ID", "IMAGE", "KERNEL", "CREATED", "SIZE", "CPUS", "MEMORY", "STATE", "IPS", "NAME")
	for _, vm := range po.VMs {
		od := vm.ObjectData.(*vmmd.VMObjectData)
		size, err := vm.Size()
		if err != nil {
			return fmt.Errorf("failed to get size for %s %q: %v", vm.Type, vm.ID, err)
		}

		imageAny, err := imgmd.LoadImageMetadata(od.ImageID)
		if err != nil {
			return fmt.Errorf("failed to load image metadata for %s %q: %v", vm.Type, vm.ID, err)
		}

		kernelAny, err := kernmd.LoadKernelMetadata(od.KernelID)
		if err != nil {
			return fmt.Errorf("failed to load kernel metadata for %s %q: %v", vm.Type, vm.ID, err)
		}

		image := imgmd.ToImageMetadata(imageAny)
		kernel := kernmd.ToKernelMetadata(kernelAny)

		// TODO: Clean up this print
		ipAddrs := []string{}
		for _, ip := range od.IPAddrs {
			ipAddrs = append(ipAddrs, ip.String())
		}
		o.Write(vm.ID, image.Name.String(), kernel.Name.String(), vm.Created, datasize.ByteSize(size).HR(), od.VCPUs,
			(datasize.ByteSize(od.Memory) * datasize.MB).HR(), od.State, strings.Join(ipAddrs, ","), vm.Name.String())
	}

	return nil
}
