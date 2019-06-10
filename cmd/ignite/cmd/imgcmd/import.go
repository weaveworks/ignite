package imgcmd

import (
	"io"

	"github.com/lithammer/dedent"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
	"github.com/weaveworks/ignite/pkg/logs"
)

// NewCmdImport imports  a new VM image
func NewCmdImport(out io.Writer) *cobra.Command {
	bo := &run.ImportOptions{}

	cmd := &cobra.Command{
		Use:   "import <source>",
		Short: "Import a new base image for VMs",
		Long: dedent.Dedent(`
			Import a new base image for VMs, takes in a Docker image as the source.
			The base image is an ext4 block device file, which contains a root filesystem.

			If the import kernel flag (-k, --import-kernel) is specified,
			/boot/vmlinux is extracted from the image and added to a new
			VM kernel object named after the flag.

			Example usage:
			    $ ignite build luxas/ubuntu-base:18.04 \
					--name my-image \
					--import-kernel my-kernel
		`),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			bo.Source = args[0]
			errutils.Check(func() error {
				var err error
				if bo.ImageNames, err = cmdutil.MatchAllImageNames(); err != nil {
					return err
				}
				return logs.PrintMachineReadableID(run.Import(bo))
			}())
		},
	}

	addImportFlags(cmd.Flags(), bo)
	return cmd
}

func addImportFlags(fs *pflag.FlagSet, bo *run.ImportOptions) {
	cmdutil.AddNameFlag(fs, &bo.Name)
	cmdutil.AddImportKernelFlags(fs, &bo.KernelName)
}
