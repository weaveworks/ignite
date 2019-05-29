package imgcmd

import (
	"github.com/spf13/cobra"
	"io"
)

// NewCmdLs lists available images
func NewCmdLs(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Short:   "List available VM base images",
		Long:    "TODO", // TODO: Long description
		Aliases: []string{"list"},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Parent().Run(cmd, args) // The parent command does this already, so just call it
		},
	}

	return cmd
}
