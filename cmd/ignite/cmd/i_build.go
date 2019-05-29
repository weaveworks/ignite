package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
)

// NewCmdImageBuild builds a new VM image
func NewCmdImageBuild(out io.Writer) *cobra.Command {
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

	addImageBuildFlags(cmd.Flags(), bo)
	return cmd
}

func addImageBuildFlags(fs *pflag.FlagSet, co *run.BuildOptions) {
	addNameFlag(fs, &co.Name)
}
