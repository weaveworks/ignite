package run

import (
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

	o.Write("VM ID", "IMAGE", "KERNEL", "CREATED", "SIZE", "CPUS", "MEMORY", "RUNNING", "IPS", "PORTS", "NAME")
	for _, vm := range po.allVMs {
		o.Write(vm.GetUID(), vm.Spec.Image.OCIClaim.Ref, vm.Spec.Kernel.OCIClaim.Ref, vm.GetCreated(),
			vm.Spec.DiskSize, vm.Spec.CPUs, vm.Spec.Memory, vm.Running(), vm.Status.IPAddresses,
			vm.Spec.Network.Ports, vm.GetName())
	}

	return nil
}
