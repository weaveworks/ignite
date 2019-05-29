package run

import (
	"fmt"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"github.com/luxas/ignite/pkg/util"
)

type StopOptions struct {
	VM *vmmd.VMMetadata
}

func Stop(so *StopOptions) error {
	// Check if the VM is running
	if !so.VM.Running() {
		return fmt.Errorf("%s is not running", so.VM.ID)
	}

	dockerArgs := []string{
		"stop",
		so.VM.ID,
	}

	// Stop the VM in docker
	if _, err := util.ExecuteCommand("docker", dockerArgs...); err != nil {
		return fmt.Errorf("failed to stop container for VM %q: %v", so.VM.ID, err)
	}

	fmt.Println(so.VM.ID)
	return nil
}
