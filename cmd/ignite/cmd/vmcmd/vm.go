package vmcmd

import (
	"io"

	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
)

// NewCmdVM handles VM-related functionality via its subcommands
func NewCmdVM(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vm",
		Short: "Manage VMs",
		Long: dedent.Dedent(`
			Groups together functionality for managing VMs.
		`),
		Aliases: []string{"vms"},
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				po, err := (&run.PsFlags{All: true}).NewPsOptions()
				if err != nil {
					return err
				}

				return run.Ps(po)
			}())
		},
	}

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
