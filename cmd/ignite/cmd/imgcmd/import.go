package imgcmd

import (
	"io"

	"github.com/lithammer/dedent"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
)

// NewCmdImport imports an image from an ext4 block device file
func NewCmdImport(out io.Writer) *cobra.Command {
	io := &run.ImportImageOptions{}

	cmd := &cobra.Command{
		Use:   "import <path>",
		Short: "Import a VM base image",
		Long: dedent.Dedent(`
			Import a new base image for VMs. This command takes in an existing ext4 block
			device file. Used in conjunction with "export" (not yet implemented).
		`), // TODO export
		Args: cobra.ExactArgs(1),
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
