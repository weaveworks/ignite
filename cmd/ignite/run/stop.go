package run

import (
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/operations"
)

type StopFlags struct {
	Kill bool
}

type StopOptions struct {
	*StopFlags
	vms []*api.VM
}

func (sf *StopFlags) NewStopOptions(vmMatches []string) (so *StopOptions, err error) {
	so = &StopOptions{StopFlags: sf}
	so.vms, err = getVMsForMatches(vmMatches)
	return
}

func Stop(so *StopOptions) error {
	for _, vm := range so.vms {
		// Stop the VM, and optionally kill it
		if err := operations.StopVM(vm, so.Kill, false); err != nil {
			return err
		}
	}

	return nil
}
