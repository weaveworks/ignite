package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
)

// NewCmdAddKernel adds a new kernel for VM use
func NewCmdAddKernel(out io.Writer) *cobra.Command {
	ao := &run.AddKernelOptions{}

	cmd := &cobra.Command{
		Use:   "addkernel [path]",
		Short: "Add an uncompressed kernel image for VM use",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ao.Source = args[0]
			errutils.Check(run.AddKernel(ao))
		},
	}

	addAddKernelFlags(cmd.Flags(), ao)
	return cmd
}

func addAddKernelFlags(fs *pflag.FlagSet, ao *run.AddKernelOptions) {
	addNameFlag(fs, &ao.Name)
}
