package run

import (
	"fmt"

	"github.com/c2h5oh/datasize"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/metadata/loader"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/util"
)

type PsFlags struct {
	All bool
}

type psOptions struct {
	*PsFlags
	allVMs []*vmmd.VMMetadata
}

func (pf *PsFlags) NewPsOptions(l *loader.ResLoader) (*psOptions, error) {
	po := &psOptions{PsFlags: pf}

	if allVMs, err := l.VMs(); err == nil {
		if po.allVMs, err = allVMs.MatchFilter(po.All); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return po, nil
}

func Ps(po *psOptions) error {
	o := util.NewOutput()
	defer o.Flush()

	o.Write("VM ID", "IMAGE", "KERNEL", "CREATED", "SIZE", "CPUS", "MEMORY", "STATE", "IPS", "PORTS", "NAME")
	for _, vm := range po.allVMs {
		od := vm.ObjectData.(*vmmd.VMObjectData)
		size, err := vm.Size()
		if err != nil {
			return fmt.Errorf("failed to get size for %s %q: %v", vm.Type, vm.ID, err)
		}

		imageAny, err := imgmd.LoadImageMetadata(od.ImageID)
		if err != nil {
			return fmt.Errorf("failed to load image metadata for %s %q: %v", vm.Type, vm.ID, err)
		}

		//kernelAny, err := kernmd.LoadKernelMetadata(od.KernelID)
		//if err != nil {
		//	return fmt.Errorf("failed to load kernel metadata for %s %q: %v", vm.Type, vm.ID, err)
		//}

		image := imgmd.ToImageMetadata(imageAny)
		//kernel := kernmd.ToKernelMetadata(kernelAny)

		// TODO: Clean up this print

		o.Write(vm.ID, image.Name.String(), "<kernel name>", vm.Created, datasize.ByteSize(size).HR(), od.VCPUs,
			od.Memory.HR(), od.State, od.IPAddrs.String(), od.PortMappings.String(), vm.Name.String())
	}

	return nil
}
