package run

import (
	"fmt"

	"github.com/weaveworks/ignite/cmd/ignite/run/runutil"

	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/util"
)

// checkRunning can be used to skip the running check, this is used by Start and Run
// as the in-container ignite takes some time to start up and update the state
type attachOptions struct {
	checkRunning bool
	vm           *vmmd.VMMetadata
}

func NewAttachOptions(l *runutil.ResLoader, vmMatch string) (*attachOptions, error) {
	ao := &attachOptions{checkRunning: true}

	if allVMs, err := l.VMs(); err == nil {
		if ao.vm, err = allVMs.MatchSingle(vmMatch); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return ao, nil
}

func Attach(ao *attachOptions) error {
	// Check if the VM is running
	if ao.checkRunning && !ao.vm.Running() {
		return fmt.Errorf("VM %q is not running", ao.vm.ID)
	}

	// Print the ID before attaching
	fmt.Println(ao.vm.ID)

	dockerArgs := []string{
		"attach",
		constants.IGNITE_PREFIX + ao.vm.ID.String(),
	}

	// Attach to the VM in Docker
	if ec, err := util.ExecForeground("docker", dockerArgs...); err != nil {
		if ec != 1 { // Docker's detach sequence (^P^Q) has an exit code of -1
			return fmt.Errorf("failed to attach to container for VM %s: %v", ao.vm.ID, err)
		}
	}

	return nil
}
