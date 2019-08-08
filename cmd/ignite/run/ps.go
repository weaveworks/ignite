package run

import (
	"fmt"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/filter"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/util"
)

type PsFlags struct {
	All bool
}

type psOptions struct {
	*PsFlags
	allVMs []*api.VM
}

func (pf *PsFlags) NewPsOptions() (po *psOptions, err error) {
	po = &psOptions{PsFlags: pf}
	po.allVMs, err = providers.Client.VMs().FindAll(filter.NewVMFilterAll("", po.All))
	return
}

func Ps(po *psOptions) error {
	o := util.NewOutput()
	defer o.Flush()

	o.Write("VM ID", "IMAGE", "KERNEL", "SIZE", "CPUS", "MEMORY", "CREATED", "STATUS", "IPS", "PORTS", "NAME")
	for _, vm := range po.allVMs {
		o.Write(vm.GetUID(), vm.Spec.Image.OCIClaim.Ref, vm.Spec.Kernel.OCIClaim.Ref,
			vm.Spec.DiskSize, vm.Spec.CPUs, vm.Spec.Memory, vm.GetCreated(), formatStatus(vm), vm.Status.IPAddresses,
			vm.Spec.Network.Ports, vm.GetName())
	}

	return nil
}

func formatStatus(vm *api.VM) (s string) {
	if vm.Running() {
		s = fmt.Sprintf("Up %s", vm.Status.StartTime)
	} else {
		s = "Stopped"
	}

	return
}
