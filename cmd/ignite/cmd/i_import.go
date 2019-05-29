package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
)

// NewCmdImageImport imports an image from an ext4 block device file
func NewCmdImageImport(out io.Writer) *cobra.Command {
	ao := &run.ImportImageOptions{}

	cmd := &cobra.Command{
		Use:   "import [path]",
		Short: "Import a VM base image",
		Long:  "TODO", // TODO: Long description
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ao.Source = args[0]
			errutils.Check(run.ImportImage(ao))
		},
	}

	addImageImportFlags(cmd.Flags(), ao)
	return cmd
}

func addImageImportFlags(fs *pflag.FlagSet, ao *run.ImportImageOptions) {
	addNameFlag(fs, &ao.Name)
}
