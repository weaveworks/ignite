package cmd

import (
	"io"

	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
)

// NewCmdStop stops a Firecracker VM.
func NewCmdStop(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop a Firecracker VM",
		Run: func(cmd *cobra.Command, args []string) {
			err := RunStop(out, cmd)
			errutils.Check(err)
		},
	}
	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

// RunStop runs when the Stop command is invoked
func RunStop(out io.Writer, cmd *cobra.Command) error {
	return nil
}
