package run

import (
	"fmt"

	"github.com/c2h5oh/datasize"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
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
		size, err := vm.Size()
		if err != nil {
			return fmt.Errorf("failed to get size for %s %q: %v", vm.Type(), vm.GetUID(), err)
		}

		imageMD, err := imgmd.LoadImage(vm.Spec.Image.ID)
		if err != nil {
			return fmt.Errorf("failed to load image metadata for %s %q: %v", vm.Type(), vm.GetUID(), err)
		}

		kernelMD, err := kernmd.LoadKernelMetadata(vm.Spec.Kernel.ID)
		if err != nil {
			return fmt.Errorf("failed to load kernel metadata for %s %q: %v", vm.Type(), vm.GetUID(), err)
		}

		image := imgmd.ToImage(imageMD)
		kernel := kernmd.ToKernelMetadata(kernelMD)

		// TODO: Clean up this print
		o.Write(vm.GetUID(), image.GetName(), kernel.GetName(), vm.Created, datasize.ByteSize(size).HR(), vm.Spec.CPUs,
			vm.Spec.Memory.HR(), vm.Status.State, vm.Status.IPAddresses, vm.Spec.Ports, vm.GetName())
	}

	return nil
}
