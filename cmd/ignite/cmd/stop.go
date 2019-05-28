package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"github.com/luxas/ignite/pkg/util"
	"github.com/spf13/cobra"
	"io"
)

type stopOptions struct {
	vm *vmmd.VMMetadata
}

// NewCmdStop stops a Firecracker VM
func NewCmdStop(out io.Writer) *cobra.Command {
	so := &stopOptions{}

	cmd := &cobra.Command{
		Use:   "stop [id]",
		Short: "Stop a running Firecracker VM",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if so.vm, err = matchSingleVM(args[0]); err != nil {
					return err
				}
				return RunStop(so)
			}())
		},
	}

	return cmd
}

func RunStop(so *stopOptions) error {
	// Check if the VM is running
	if !so.vm.Running() {
		return fmt.Errorf("%s is not running", so.vm.ID)
	}

	dockerArgs := []string{
		"stop",
		so.vm.ID,
	}

	// Stop the VM in docker
	if _, err := util.ExecuteCommand("docker", dockerArgs...); err != nil {
		return fmt.Errorf("failed to stop container for VM %q: %v", so.vm.ID, err)
	}

	fmt.Println(so.vm.ID)
	return nil
}
