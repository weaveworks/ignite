package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
)

// NewCmdBuild builds a new VM base image
func NewCmdBuild(out io.Writer) *cobra.Command {
	bo := &run.BuildOptions{}

	cmd := &cobra.Command{
		Use:   "build [source]",
		Short: "Build a Firecracker VM base image",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			bo.Source = args[0]
			errutils.Check(run.Build(bo))
		},
	}

	addBuildFlags(cmd.Flags(), bo)
	return cmd
}

func addBuildFlags(fs *pflag.FlagSet, co *run.BuildOptions) {
	addNameFlag(fs, &co.Name)
}
