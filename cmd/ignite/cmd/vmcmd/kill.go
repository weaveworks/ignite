package vmcmd

import (
	"github.com/luxas/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdKill kills a VM
func NewCmdKill(out io.Writer) *cobra.Command {
	so := &run.StopOptions{Kill: true}

	cmd := &cobra.Command{
		Use:   "kill [vm]...",
		Short: "Kill running VMs",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if so.VMs, err = cmdutil.MatchSingleVMs(args); err != nil {
					return err
				}
				return run.Stop(so)
			}())
		},
	}

	return cmd
}
