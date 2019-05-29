package vmcmd

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

	cmd.AddCommand(NewCmdAttach(out))
	cmd.AddCommand(NewCmdCreate(out))
	cmd.AddCommand(NewCmdKill(out))
	cmd.AddCommand(NewCmdLogs(out))
	cmd.AddCommand(NewCmdPs(out))
	cmd.AddCommand(NewCmdRm(out))
	cmd.AddCommand(NewCmdRun(out))
	cmd.AddCommand(NewCmdStart(out))
	cmd.AddCommand(NewCmdStop(out))
	return cmd
}
