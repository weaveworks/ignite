package imgcmd

import (
	"github.com/luxas/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
)

// NewCmdBuild builds a new VM image
func NewCmdBuild(out io.Writer) *cobra.Command {
	bo := &run.BuildOptions{}

	cmd := &cobra.Command{
		Use:   "build [source]",
		Short: "Build a VM base image",
		Long:  "TODO", // TODO: Long description
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			bo.Source = args[0]
			errutils.Check(run.Build(bo))
		},
	}

	addBuildFlags(cmd.Flags(), bo)
	return cmd
}

func addBuildFlags(fs *pflag.FlagSet, bo *run.BuildOptions) {
	cmdutil.AddNameFlag(fs, &bo.Name)
	cmdutil.AddImportKernelFlags(fs, &bo.KernelName)
}
