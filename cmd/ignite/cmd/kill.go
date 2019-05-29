package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdStop kills a Firecracker VM
func NewCmdKill(out io.Writer) *cobra.Command {
	ko := &run.KillOptions{}

	cmd := &cobra.Command{
		Use:   "kill [id]",
		Short: "Kill a running Firecracker VM",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if ko.VM, err = matchSingleVM(args[0]); err != nil {
					return err
				}
				return run.Kill(ko)
			}())
		},
	}

	return cmd
}
