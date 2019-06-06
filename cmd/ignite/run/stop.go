package run

import (
	"fmt"

	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"github.com/luxas/ignite/pkg/util"
)

var (
	stopArgs = []string{"stop"}
	killArgs = []string{"kill", "-s", "SIGQUIT"}
)

type StopOptions struct {
	VMs  []*vmmd.VMMetadata
	Kill bool
}

func Stop(so *StopOptions) error {
	for _, vm := range so.VMs {
		// Check if the VM is running
		if !vm.Running() {
			return fmt.Errorf("%s is not running", vm.ID)
		}

		dockerArgs := stopArgs

		// Change to kill arguments if requested
		if so.Kill {
			dockerArgs = killArgs
		}

		dockerArgs = append(dockerArgs, constants.IGNITE_PREFIX+vm.ID)

		// Stop/Kill the VM in docker
		if _, err := util.ExecuteCommand("docker", dockerArgs...); err != nil {
			return fmt.Errorf("failed to stop container for VM %q: %v", vm.ID, err)
		}

		fmt.Println(vm.ID)
	}

	return nil
}
