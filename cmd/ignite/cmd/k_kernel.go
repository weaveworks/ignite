package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdKernel handles kernel-related functionality via its subcommands
// This command by itself lists available kernels
func NewCmdKernel(out io.Writer) *cobra.Command {
	ko := &run.KernelsOptions{}

	cmd := &cobra.Command{
		Use:     "kernel",
		Short:   "Manage VM kernels",
		Long:    "TODO", // TODO: Long description
		Aliases: []string{"kernels"},
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

	cmd.AddCommand(NewCmdKernelImport(out))
	cmd.AddCommand(NewCmdKernelLs(out))
	cmd.AddCommand(NewCmdKernelRm(out))
	return cmd
}
