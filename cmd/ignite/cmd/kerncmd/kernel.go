package kerncmd

import (
	"io"

	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
	"github.com/spf13/cobra"
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
				if ko.Kernels, err = cmdutil.MatchAllKernels(); err != nil {
					return err
				}
				return run.Kernels(ko)
			}())
		},
	}

	cmd.AddCommand(NewCmdImport(out))
	cmd.AddCommand(NewCmdLs(out))
	cmd.AddCommand(NewCmdRm(out))
	return cmd
}
