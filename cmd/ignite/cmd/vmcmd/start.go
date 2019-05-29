package vmcmd

import (
	"github.com/luxas/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
)

// NewCmdStart starts a VM
func NewCmdStart(out io.Writer) *cobra.Command {
	so := &run.StartOptions{}

	cmd := &cobra.Command{
		Use:   "start [vm]",
		Short: "Start a VM",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if so.VM, err = cmdutil.MatchSingleVM(args[0]); err != nil {
					return err
				}
				return run.Start(so)
			}())
		},
	}

	addStartFlags(cmd.Flags(), so)
	return cmd
}

func addStartFlags(fs *pflag.FlagSet, so *run.StartOptions) {
	cmdutil.AddInteractiveFlag(fs, &so.Interactive)
}
