package run

import (
	"fmt"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"github.com/luxas/ignite/pkg/util"
)

var (
	stopArgs = []string{"stop"}
	killArgs = []string{"kill", "-s", "SIGQUIT"}
)

type StopOptions struct {
	VM   *vmmd.VMMetadata
	Kill bool
}

func Stop(so *StopOptions) error {
	// Check if the VM is running
	if !so.VM.Running() {
		return fmt.Errorf("%s is not running", so.VM.ID)
	}

	dockerArgs := stopArgs

	// Change to kill arguments if requested
	if so.Kill {
		dockerArgs = killArgs
	}

	dockerArgs = append(dockerArgs, so.VM.ID)

	// Stop/Kill the VM in docker
	if _, err := util.ExecuteCommand("docker", dockerArgs...); err != nil {
		return fmt.Errorf("failed to stop container for VM %q: %v", so.VM.ID, err)
	}

	fmt.Println(so.VM.ID)
	return nil
}
