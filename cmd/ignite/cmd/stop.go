package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdStop stops a Firecracker VM
func NewCmdStop(out io.Writer) *cobra.Command {
	so := &run.StopOptions{}

	cmd := &cobra.Command{
		Use:   "stop [id]",
		Short: "Stop a running Firecracker VM",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if so.VM, err = matchSingleVM(args[0]); err != nil {
					return err
				}
				return run.Stop(so)
			}())
		},
	}

	return cmd
}
