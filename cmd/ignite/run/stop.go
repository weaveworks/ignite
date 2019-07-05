package run

import (
	"fmt"

	"github.com/weaveworks/ignite/pkg/metadata/loader"
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

func (sf *StopFlags) NewStopOptions(l *loader.ResLoader, vmMatches []string) (*stopOptions, error) {
	so := &stopOptions{StopFlags: sf}

	if allVMs, err := l.VMs(); err == nil {
		if so.vms, err = allVMs.MatchMultiple(vmMatches); err != nil {
			return nil, err
		}
	} else {
		return nil, err
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
		if err := operations.StopVM(vm.VM, so.Kill, false); err != nil {
			return err
		}
	}
	return nil
}
