package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
)

// NewCmdKernelImport imports a kernel for VM use
func NewCmdKernelImport(out io.Writer) *cobra.Command {
	ao := &run.ImportKernelOptions{}

	cmd := &cobra.Command{
		Use:   "import [path]",
		Short: "Import an uncompressed kernel image for VM use",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ao.Source = args[0]
			errutils.Check(run.ImportKernel(ao))
		},
	}

	addKernelImportFlags(cmd.Flags(), ao)
	return cmd
}

func addKernelImportFlags(fs *pflag.FlagSet, ao *run.ImportKernelOptions) {
	addNameFlag(fs, &ao.Name)
}
