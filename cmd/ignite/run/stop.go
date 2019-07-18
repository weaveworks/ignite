package run

import (
	"fmt"

	"github.com/weaveworks/ignite/pkg/runtime/docker"

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

func (sf *StopFlags) NewStopOptions(vmMatches []string) (so *stopOptions, err error) {
	so = &stopOptions{StopFlags: sf}
	so.vms, err = getVMsForMatches(vmMatches)
	return
}

func Stop(so *stopOptions) error {
	// Get the Docker client
	dc, err := docker.GetDockerClient()
	if err != nil {
		return err
	}

	for _, vm := range so.vms {
		// Check if the VM is running
		if !vm.Running() {
			return fmt.Errorf("VM %q is not running", vm.GetUID())
		}

		// Stop the VM, and optionally kill it
		if err := operations.StopVM(dc, vm, so.Kill, false); err != nil {
			return err
		}
	}

	return nil
}
