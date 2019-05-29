package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
)

// NewCmdAddImage imports an image for VM use
func NewCmdAddImage(out io.Writer) *cobra.Command {
	ao := &run.AddImageOptions{}

	cmd := &cobra.Command{
		Use:   "addimage [path]",
		Short: "Import an existing VM base image",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ao.Source = args[0]
			errutils.Check(run.AddImage(ao))
		},
	}

	addAddImageFlags(cmd.Flags(), ao)
	return cmd
}

func addAddImageFlags(fs *pflag.FlagSet, ao *run.AddImageOptions) {
	addNameFlag(fs, &ao.Name)
}
