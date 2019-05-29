package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdKernelRm removes a kernel
// TODO: Support removing multiple kernels at once
func NewCmdKernelRm(out io.Writer) *cobra.Command {
	ro := &run.RmkOptions{}

	cmd := &cobra.Command{
		Use:   "rm [kernel]",
		Short: "Remove a kernel",
		Long:  "TODO", // TODO: Long description
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
