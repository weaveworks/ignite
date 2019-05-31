package imgcmd

import (
	"github.com/luxas/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
)

// NewCmdImport imports an image from an ext4 block device file
func NewCmdImport(out io.Writer) *cobra.Command {
	io := &run.ImportImageOptions{}

	cmd := &cobra.Command{
		Use:   "import [path]",
		Short: "Import a VM base image",
		Long:  "TODO", // TODO: Long description
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			io.Source = args[0]
			errutils.Check(run.ImportImage(io))
		},
	}

	addImportFlags(cmd.Flags(), io)
	return cmd
}

func addImportFlags(fs *pflag.FlagSet, io *run.ImportImageOptions) {
	cmdutil.AddNameFlag(fs, &io.Name)
	cmdutil.AddImportKernelFlags(fs, &io.ImportKernel, &io.KernelName)
}
