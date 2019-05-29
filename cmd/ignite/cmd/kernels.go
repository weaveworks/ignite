package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdKernels lists the available kernels
func NewCmdKernels(out io.Writer) *cobra.Command {
	ko := &run.KernelsOptions{}

	cmd := &cobra.Command{
		Use:   "kernels",
		Short: "List available kernels",
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if ko.Kernels, err = matchAllKernels(); err != nil {
					return err
				}
				return run.Kernels(ko)
			}())
		},
	}

	return cmd
}
