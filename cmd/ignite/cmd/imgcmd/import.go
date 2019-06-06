package imgcmd

import (
	"io"

	"github.com/luxas/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
			errutils.Check(func() error {
				var err error
				if io.ImageNames, err = cmdutil.MatchAllImageNames(); err != nil {
					return err
				}
				return run.ImportImage(io)
			}())
		},
	}

	addImportFlags(cmd.Flags(), io)
	return cmd
}

func addImportFlags(fs *pflag.FlagSet, io *run.ImportImageOptions) {
	cmdutil.AddNameFlag(fs, &io.Name)
	cmdutil.AddImportKernelFlags(fs, &io.KernelName)
}
