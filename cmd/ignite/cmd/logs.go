package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"github.com/luxas/ignite/pkg/util"
	"io"

	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
)

type logsOptions struct {
	vm *vmmd.VMMetadata
}

// NewCmdLogs gets the logs for a Firecracker VM
func NewCmdLogs(out io.Writer) *cobra.Command {
	lo := &logsOptions{}

	cmd := &cobra.Command{
		Use:   "logs [id]",
		Short: "Gets the logs for a Firecracker VM",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if lo.vm, err = matchSingleVM(args[0]); err != nil {
					return err
				}
				return RunLogs(lo)
			}())
		},
	}

	return cmd
}

func RunLogs(lo *logsOptions) error {
	// Check if the VM is running
	if !lo.vm.Running() {
		return fmt.Errorf("%s is not running", lo.vm.ID)
	}

	dockerArgs := []string{
		"logs",
		lo.vm.ID,
	}

	// Fetch the VM logs from docker
	output, err := util.ExecuteCommand("docker", dockerArgs...)
	if err != nil {
		return fmt.Errorf("failed to get logs for VM %q: %v", lo.vm.ID, err)
	}

	// Print the ID and the VM logs
	fmt.Println(lo.vm.ID)
	fmt.Println(output)
	return nil
}
