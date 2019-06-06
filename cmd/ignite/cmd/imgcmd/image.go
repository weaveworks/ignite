package imgcmd

import (
	"io"

	"github.com/luxas/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
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
				if io.Images, err = cmdutil.MatchAllImages(); err != nil {
					return err
				}
				return run.Images(io)
			}())
		},
	}

	cmd.AddCommand(NewCmdBuild(out))
	cmd.AddCommand(NewCmdImport(out))
	cmd.AddCommand(NewCmdLs(out))
	cmd.AddCommand(NewCmdRm(out))
	return cmd
}
