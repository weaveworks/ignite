package cmd

import (
	"github.com/spf13/cobra"
	"io"
)

func NewCmdHang(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "hang",
		Short:  "Causes Ignite to hang indefinitely. Used for testing purposes.",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			select {}
		},
	}

	return cmd
}
