package run

import (
	"fmt"

	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/util"
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
			return fmt.Errorf("VM %q is not running", vm.ID)
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
