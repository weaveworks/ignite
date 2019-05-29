package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdVMKill kills a VM
func NewCmdVMKill(out io.Writer) *cobra.Command {
	ko := &run.KillOptions{}

	cmd := &cobra.Command{
		Use:   "kill [vm]",
		Short: "Kill a running VM",
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
