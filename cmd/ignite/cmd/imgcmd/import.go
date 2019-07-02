package imgcmd

import (
	"io"

	"github.com/weaveworks/ignite/pkg/metadata/loader"

	"github.com/lithammer/dedent"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
)

// NewCmdImport imports  a new VM image
func NewCmdImport(out io.Writer) *cobra.Command {
	ifs := &run.ImportFlags{}

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
			errutils.Check(func() error {
				i, err := ifs.NewImportOptions(loader.NewResLoader(), args[0])
				if err != nil {
					return err
				}
				return run.Import(i)
			}())
		},
	}

	addImportFlags(cmd.Flags(), ifs)
	return cmd
}

func addImportFlags(fs *pflag.FlagSet, ifs *run.ImportFlags) {
	cmdutil.AddNameFlag(fs, &ifs.Name)
}
