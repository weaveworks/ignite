package run

import (
	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/filter"
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
	po.allVMs, err = client.VMs().FindAll(filter.NewVMFilterAll("", po.All))
	return
}

func Ps(po *psOptions) error {
	o := util.NewOutput()
	defer o.Flush()

	o.Write("VM ID", "IMAGE", "KERNEL", "CREATED", "SIZE", "CPUS", "MEMORY", "STATE", "IPS", "PORTS", "NAME")
	for _, vm := range po.allVMs {
		o.Write(vm.GetUID(), vm.Spec.Image.OCIClaim.Ref.String(), vm.Spec.Kernel.OCIClaim.Ref.String(), vm.GetCreated(),
			vm.Spec.DiskSize.String(), vm.Spec.CPUs, vm.Spec.Memory.String(), vm.Status.State, vm.Status.IPAddresses,
			vm.Spec.Ports, vm.GetName())
	}

	return nil
}
