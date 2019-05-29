package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdAddImage imports an image for VM use
func NewCmdAddImage(out io.Writer) *cobra.Command {
	ao := &run.ImportImageOptions{}

	cmd := &cobra.Command{
		Use:   "addimage [path]",
		Short: "Import an existing VM base image",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ao.Source = args[0]
			errutils.Check(run.ImportImage(ao))
		},
	}

	addImageImportFlags(cmd.Flags(), ao)
	return cmd
}
