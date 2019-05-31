package kerncmd

import (
	"github.com/luxas/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
)

// NewCmdImport imports a kernel for VM use
func NewCmdImport(out io.Writer) *cobra.Command {
	io := &run.ImportKernelOptions{}

	cmd := &cobra.Command{
		Use:   "import [path]",
		Short: "Import an uncompressed kernel image for VM use",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			io.Source = args[0]
			errutils.Check(run.ImportKernel(io))
		},
	}

	addImportFlags(cmd.Flags(), io)
	return cmd
}

func addImportFlags(fs *pflag.FlagSet, io *run.ImportKernelOptions) {
	cmdutil.AddNameFlag(fs, &io.Name)
}
