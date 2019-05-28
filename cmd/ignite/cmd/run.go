package cmd

import (
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdRun creates, starts (and attaches to) a Firecracker VM
func NewCmdRun(out io.Writer) *cobra.Command {
	co := &createOptions{}

	cmd := &cobra.Command{
		Use:   "run [image] [kernel] [name]",
		Short: "Create and start a new Firecracker VM",
		Args:  cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			err := RunRun(out, cmd, args[0], args[1], args[2], co)
			errutils.Check(err)
		},
	}

	cmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Immediately attach to created and started VM")
	cmd.Flags().Int64Var(&co.cpus, "cpus", constants.VM_DEFAULT_CPUS, "VM vCPU count, 1 or even numbers between 1 and 32")
	cmd.Flags().Int64Var(&co.memory, "memory", constants.VM_DEFAULT_MEMORY, "VM RAM in MiB")
	return cmd
}

func RunRun(out io.Writer, cmd *cobra.Command, imageMatch, kernelMatch, name string, co *createOptions) error {
	if err := RunCreate(out, cmd, imageMatch, kernelMatch, name, co, true); err != nil {
		return err
	}

	return nil
}
