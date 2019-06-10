package vmcmd

import (
	"io"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
)

// NewCmdVM handles VM-related functionality via its subcommands
func NewCmdVM(out io.Writer) *cobra.Command {
	po := &run.PsOptions{}

	cmd := &cobra.Command{
		Use:   "vm",
		Short: "Manage VMs",
		Long: dedent.Dedent(`
			Groups together functionality for managing VMs.
		`),
		Aliases: []string{"vms"},
		Run: func(cmd *cobra.Command, args []string) {
			po.All = true
			errutils.Check(ExecutePs(po))
		},
	}

	addPsFlags(cmd.Flags(), po)

	cmd.AddCommand(NewCmdAttach(out))
	cmd.AddCommand(NewCmdCreate(out))
	cmd.AddCommand(NewCmdKill(out))
	cmd.AddCommand(NewCmdLogs(out))
	cmd.AddCommand(NewCmdPs(out))
	cmd.AddCommand(NewCmdRm(out))
	cmd.AddCommand(NewCmdRun(out))
	cmd.AddCommand(NewCmdSSH(out))
	cmd.AddCommand(NewCmdStart(out))
	cmd.AddCommand(NewCmdStop(out))
	return cmd
}
