package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"github.com/luxas/ignite/pkg/util"
	"github.com/spf13/cobra"
	"io"
)

type killOptions struct {
	vm *vmmd.VMMetadata
}

// NewCmdStop kills a Firecracker VM
func NewCmdKill(out io.Writer) *cobra.Command {
	ko := &killOptions{}

	cmd := &cobra.Command{
		Use:   "kill [id]",
		Short: "Kill a running Firecracker VM",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if ko.vm, err = matchSingleVM(args[0]); err != nil {
					return err
				}
				return RunKill(ko)
			}())
		},
	}

	return cmd
}

func RunKill(ko *killOptions) error {
	// Check if the VM is running
	if !ko.vm.Running() {
		return fmt.Errorf("%s is not running", ko.vm.ID)
	}

	dockerArgs := []string{
		"kill",
		"-s",
		"SIGQUIT",
		ko.vm.ID,
	}

	// Kill the VM in docker
	if _, err := util.ExecuteCommand("docker", dockerArgs...); err != nil {
		return fmt.Errorf("failed to kill container for VM %q: %v", ko.vm.ID, err)
	}

	fmt.Println(ko.vm.ID)
	return nil
}
