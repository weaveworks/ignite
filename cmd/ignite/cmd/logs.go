package cmd

import (
	"io"

	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
)

// NewCmdLogs gets the logs for a Firecracker VM.
func NewCmdLogs(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Gets the logs for a Firecracker VM",
		Run: func(cmd *cobra.Command, args []string) {
			err := RunLogs(out, cmd)
			errutils.Check(err)
		},
	}
	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

// RunLogs runs when the Logs command is invoked
func RunLogs(out io.Writer, cmd *cobra.Command) error {
	return nil
}
