package run

import (
	"fmt"

	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/filter"

	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/util"
)

// checkRunning can be used to skip the running check, this is used by Start and Run
// as the in-container ignite takes some time to start up and update the state
type attachOptions struct {
	checkRunning bool
	vm           *vmmd.VM
}

func NewAttachOptions(vmMatch string) (*attachOptions, error) {
	ao := &attachOptions{checkRunning: true}

	if vm, err := client.VMs().Find(filter.NewIDNameFilter(vmMatch)); err == nil {
		ao.vm = &vmmd.VM{vm}
	} else {
		return nil, err
	}

	return ao, nil
}

func Attach(ao *attachOptions) error {
	// Check if the VM is running
	if ao.checkRunning && !ao.vm.Running() {
		return fmt.Errorf("VM %q is not running", ao.vm.GetUID())
	}

	// Print the ID before attaching
	fmt.Println(ao.vm.GetUID())

	dockerArgs := []string{
		"attach",
		constants.IGNITE_PREFIX + ao.vm.GetUID().String(),
	}

	// Attach to the VM in Docker
	if ec, err := util.ExecForeground("docker", dockerArgs...); err != nil {
		if ec != 1 { // Docker's detach sequence (^P^Q) has an exit code of -1
			return fmt.Errorf("failed to attach to container for VM %s: %v", ao.vm.GetUID(), err)
		}
	}

	return nil
}
