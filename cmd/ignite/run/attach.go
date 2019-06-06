package run

import (
	"fmt"

	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/util"
)

type AttachOptions struct {
	VM           *vmmd.VMMetadata
	CheckRunning bool
}

func Attach(ao *AttachOptions) error {
	// Check if the VM is running
	if ao.CheckRunning && !ao.VM.Running() {
		return fmt.Errorf("%s is not running", ao.VM.ID)
	}

	// Print the ID before attaching
	fmt.Println(ao.VM.ID)

	dockerArgs := []string{
		"attach",
		constants.IGNITE_PREFIX + ao.VM.ID,
	}

	// Attach to the VM in Docker
	if ec, err := util.ExecForeground("docker", dockerArgs...); err != nil {
		if ec != 1 { // Docker's detach sequence (^P^Q) has an exit code of -1
			return fmt.Errorf("failed to attach to container for VM %s: %v", ao.VM.ID, err)
		}
	}

	return nil
}
