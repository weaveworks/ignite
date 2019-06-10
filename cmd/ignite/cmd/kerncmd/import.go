package kerncmd

import (
	"github.com/weaveworks/ignite/cmd/ignite/run/runutil"
	"io"

	"github.com/lithammer/dedent"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
)

// NewCmdImport imports a kernel for VM use
func NewCmdImport(out io.Writer) *cobra.Command {
	io := &run.ImportKernelOptions{}

	cmd := &cobra.Command{
		Use:   "import <path>",
		Short: "Import an uncompressed kernel image for VM use",
		Args:  cobra.ExactArgs(1),
		Long: dedent.Dedent(`
			Import a new kernel for VMs. This command takes in an existing uncompressed
			kernel (vmlinux) file. Used in conjunction with "export" (not yet implemented).
		`), // TODO export
		Run: func(cmd *cobra.Command, args []string) {
			io.Source = args[0]
			errutils.Check(func() error {
				var err error
				if io.KernelNames, err = runutil.MatchAllKernelNames(); err != nil {
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
