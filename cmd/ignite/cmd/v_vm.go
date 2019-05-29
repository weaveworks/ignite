package cmd

import (
	"github.com/spf13/cobra"
	"io"
)

// NewCmdVM handles VM-related functionality via its subcommands
func NewCmdVM(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vm",
		Short:   "Manage VMs",
		Long:    "TODO", // TODO: Long descriptions for this and subcommands
		Aliases: []string{"vms"},
	}

	cmd.AddCommand(NewCmdVMAttach(out))
	cmd.AddCommand(NewCmdVMCreate(out))
	cmd.AddCommand(NewCmdVMKill(out))
	cmd.AddCommand(NewCmdVMLogs(out))
	cmd.AddCommand(NewCmdVMPs(out))
	cmd.AddCommand(NewCmdVMRm(out))
	cmd.AddCommand(NewCmdVMRun(out))
	cmd.AddCommand(NewCmdVMStart(out))
	cmd.AddCommand(NewCmdVMStop(out))
	return cmd
}
