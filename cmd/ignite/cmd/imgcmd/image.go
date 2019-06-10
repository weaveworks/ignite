package imgcmd

import (
	"io"

	"github.com/lithammer/dedent"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
)

// NewCmdImage handles image-related functionality via its subcommands
// This command by itself lists available images
func NewCmdImage(out io.Writer) *cobra.Command {
	io := &run.ImagesOptions{}

	cmd := &cobra.Command{
		Use:   "image",
		Short: "Manage VM base images",
		Long: dedent.Dedent(`
			Groups together functionality for managing VM base images.
			Calling this command alone lists all available images.
		`),
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

	cmd.AddCommand(NewCmdImport(out))
	cmd.AddCommand(NewCmdLs(out))
	cmd.AddCommand(NewCmdRm(out))
	return cmd
}
