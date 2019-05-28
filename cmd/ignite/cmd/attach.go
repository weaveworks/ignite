package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"github.com/luxas/ignite/pkg/util"
	"github.com/spf13/cobra"
	"io"
)

type attachOptions struct {
	vm           *vmmd.VMMetadata
	checkRunning bool
}

// NewCmdAttach attaches to a running Firecracker VM
func NewCmdAttach(out io.Writer) *cobra.Command {
	// checkRunning can be used to skip the running check, this is used by Start and Run
	// as the in-container ignite takes some time to start up and update the state
	ao := &attachOptions{checkRunning: true}

	cmd := &cobra.Command{
		Use:   "attach [vm]",
		Short: "Attach to a running Firecracker VM",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if ao.vm, err = matchSingleVM(args[0]); err != nil {
					return err
				}
				return RunAttach(ao)
			}())
		},
	}

	return cmd
}

func RunAttach(ao *attachOptions) error {
	// Check if the VM is running
	if ao.checkRunning && !ao.vm.Running() {
		return fmt.Errorf("%s is not running", ao.vm.ID)
	}

	// Print the ID before attaching
	fmt.Println(ao.vm.ID)

	dockerArgs := []string{
		"attach",
		ao.vm.ID,
	}

	// Attach to the VM in Docker
	if ec, err := util.ExecForeground("docker", dockerArgs...); err != nil {
		if ec != 1 { // Docker's detach sequence (^P^Q) has an exit code of -1
			return fmt.Errorf("failed to attach to container for VM %s: %v", ao.vm.ID, err)
		}
	}

	return nil
}
