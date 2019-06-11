package run

import (
	"fmt"

	"github.com/weaveworks/ignite/cmd/ignite/run/runutil"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/util"
)

var (
	stopArgs = []string{"stop"}
	killArgs = []string{"kill", "-s", "SIGQUIT"}
)

type StopFlags struct {
	Kill bool
}

type stopOptions struct {
	*StopFlags
	vms    []*vmmd.VMMetadata
	silent bool
}

func (sf *StopFlags) NewStopOptions(l *runutil.ResLoader, vmMatches []string) (*stopOptions, error) {
	so := &stopOptions{}

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
			return fmt.Errorf("VM %q is not running", vm.ID)
		}

		dockerArgs := stopArgs

		// Change to kill arguments if requested
		if so.Kill {
			dockerArgs = killArgs
		}

		dockerArgs = append(dockerArgs, constants.IGNITE_PREFIX+vm.ID.String())

		// Stop/Kill the VM in docker
		if _, err := util.ExecuteCommand("docker", dockerArgs...); err != nil {
			return fmt.Errorf("failed to stop container for VM %q: %v", vm.ID, err)
		}

		if so.silent {
			continue
		}

		fmt.Println(vm.ID)
	}
	return nil
}
