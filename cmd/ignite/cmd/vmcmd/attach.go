package vmcmd

import (
	"github.com/luxas/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdAttach attaches to a running VM
func NewCmdAttach(out io.Writer) *cobra.Command {
	// checkRunning can be used to skip the running check, this is used by Start and Run
	// as the in-container ignite takes some time to start up and update the state
	ao := &run.AttachOptions{CheckRunning: true}

	cmd := &cobra.Command{
		Use:   "attach [vm]",
		Short: "Attach to a running VM",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if ao.VM, err = cmdutil.MatchSingleVM(args[0]); err != nil {
					return err
				}
				return run.Attach(ao)
			}())
		},
	}

	return cmd
}
