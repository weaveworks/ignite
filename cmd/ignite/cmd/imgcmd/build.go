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

// NewCmdBuild builds a new VM image
func NewCmdBuild(out io.Writer) *cobra.Command {
	bo := &run.BuildOptions{}

	cmd := &cobra.Command{
		Use:   "build [source]",
		Short: "Build a new base image for VMs",
		Long: dedent.Dedent(`
			Build a new base image for VMs. The base image is an ext4
			block device file, which contains a root filesystem.
			
			"build" can take in either a tarfile or a Docker image as the source.
			The Docker image needs to exist on the host system (pulled locally).

			If the import kernel flag (-k, --import-kernel) is specified,
			/boot/vmlinux is extracted from the image and added to a new
			VM kernel object named after the flag.

			Example usage:
				$ ignite build my-image.tar
			    $ ignite build luxas/ubuntu-base:18.04 \
					--name my-image \
					--import-kernel my-kernel
		`),
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			bo.Source = args[0]
			errutils.Check(func() error {
				var err error
				if bo.ImageNames, err = cmdutil.MatchAllImageNames(); err != nil {
					return err
				}
				return run.Build(bo)
			}())
		},
	}

	addBuildFlags(cmd.Flags(), bo)
	return cmd
}

func addBuildFlags(fs *pflag.FlagSet, bo *run.BuildOptions) {
	cmdutil.AddNameFlag(fs, &bo.Name)
	cmdutil.AddImportKernelFlags(fs, &bo.KernelName)
}
