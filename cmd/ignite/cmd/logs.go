package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/run"
	"io"

	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
)

// NewCmdLogs gets the logs for a Firecracker VM
func NewCmdLogs(out io.Writer) *cobra.Command {
	lo := &run.LogsOptions{}

	cmd := &cobra.Command{
		Use:   "logs [id]",
		Short: "Gets the logs for a Firecracker VM",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if lo.VM, err = matchSingleVM(args[0]); err != nil {
					return err
				}
				return run.Logs(lo)
			}())
		},
	}

	return cmd
}
