package cmd

import (
	"github.com/spf13/cobra"
	"io"
)

// NewCmdKernelLs lists available kernels
func NewCmdKernelLs(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Short:   "List available VM kernels",
		Long:    "TODO", // TODO: Long description
		Aliases: []string{"list"},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Parent().Run(cmd, args) // The parent command does this already, so just call it
		},
	}

	return cmd
}
