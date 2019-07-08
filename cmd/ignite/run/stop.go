package run

import (
	"fmt"

	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/filter"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/operations"
)

type StopFlags struct {
	Kill bool
}

type stopOptions struct {
	*StopFlags
	vms []*vmmd.VM
}

func (sf *StopFlags) NewStopOptions(vmMatches []string) (*stopOptions, error) {
	so := &stopOptions{StopFlags: sf}

	for _, match := range vmMatches {
		if vm, err := client.VMs().Find(filter.NewIDNameFilter(match)); err == nil {
			so.vms = append(so.vms, &vmmd.VM{vm})
		} else {
			return nil, err
		}
	}

	return so, nil
}

func Stop(so *stopOptions) error {
	for _, vm := range so.vms {
		// Check if the VM is running
		if !vm.Running() {
			return fmt.Errorf("VM %q is not running", vm.GetUID())
		}

		// Stop the VM, and optionally kill it
		if err := operations.StopVM(vm, so.Kill, false); err != nil {
			return err
		}
	}

	return nil
}
