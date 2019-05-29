package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdRmk removes the given kernel
func NewCmdRmk(out io.Writer) *cobra.Command {
	ro := &run.RmkOptions{}

	cmd := &cobra.Command{
		Use:   "rmk [id]",
		Short: "Remove a kernel",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if ro.Kernel, err = matchSingleKernel(args[0]); err != nil {
					return err
				}
				return run.Rmk(ro)
			}())
		},
	}

	return cmd
}
