package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdImage handles image-related functionality via its subcommands
// This command by itself lists available images
func NewCmdImage(out io.Writer) *cobra.Command {
	io := &run.ImagesOptions{}

	cmd := &cobra.Command{
		Use:     "image",
		Short:   "Manage VM base images",
		Long:    "TODO", // TODO: Long description
		Aliases: []string{"images"},
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if io.Images, err = matchAllImages(); err != nil {
					return err
				}
				return run.Images(io)
			}())
		},
	}

	cmd.AddCommand(NewCmdImageBuild(out))
	cmd.AddCommand(NewCmdImageImport(out))
	cmd.AddCommand(NewCmdImageLs(out))
	cmd.AddCommand(NewCmdImageRm(out))
	return cmd
}
