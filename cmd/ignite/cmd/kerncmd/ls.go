package kerncmd

import (
	"io"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
)

// NewCmdLs lists available kernels
func NewCmdLs(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List available VM kernels",
		Long: dedent.Dedent(`
			List all available VM kernels. Outputs the same as the parent command.
		`),
		Aliases: []string{"list"},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Parent().Run(cmd, args) // The parent command does this already, so just call it
		},
	}

	return cmd
}
