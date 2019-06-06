package kerncmd

import (
	"io"

	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
			errutils.Check(func() error {
				var err error
				if io.KernelNames, err = cmdutil.MatchAllKernelNames(); err != nil {
					return err
				}
				return run.ImportKernel(io)
			}())
		},
	}

	addImportFlags(cmd.Flags(), io)
	return cmd
}

func addImportFlags(fs *pflag.FlagSet, io *run.ImportKernelOptions) {
	cmdutil.AddNameFlag(fs, &io.Name)
}
