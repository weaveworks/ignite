package run

import (
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/config"
	"github.com/weaveworks/ignite/pkg/operations"
)

type StopFlags struct {
	Kill bool
}

type stopOptions struct {
	*StopFlags
	vms []*api.VM
}

func (sf *StopFlags) NewStopOptions(vmMatches []string) (so *stopOptions, err error) {
	so = &stopOptions{StopFlags: sf}
	so.vms, err = getVMsForMatches(vmMatches)
	return
}

func Stop(so *stopOptions) error {
	for _, vm := range so.vms {
		// Set the runtime and network-plugin providers from the VM status.
		if err := config.SetAndPopulateProviders(vm.Status.Runtime.Name, vm.Status.Network.Plugin); err != nil {
			return err
		}

		// Stop the VM, and optionally kill it
		if err := operations.StopVM(vm, so.Kill, false); err != nil {
			return err
		}
	}

	return nil
}
