package imgcmd

import (
	"io"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
)

// NewCmdImage handles image-related functionality via its subcommands
// This command by itself lists available images
func NewCmdImage(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "Manage base images for VMs",
		Long: dedent.Dedent(`
			Groups together functionality for managing VM base images.
			Calling this command alone lists all available images.
		`),
		Aliases: []string{"images"},
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(func() error {
				i, err := run.NewImagesOptions()
				if err != nil {
					return err
				}

				return run.Images(i)
			}())
		},
	}

	cmd.AddCommand(NewCmdImport(out))
	cmd.AddCommand(NewCmdLs(out))
	cmd.AddCommand(NewCmdRm(out))
	return cmd
}
