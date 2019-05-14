package cmd

import (
	"io"

	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
)

// NewCmdExec execs a command in a Firecracker VM.
func NewCmdExec(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec",
		Short: "Execute a command in a Firecracker VM",
		Run: func(cmd *cobra.Command, args []string) {
			err := RunExec(out, cmd)
			errutils.Check(err)
		},
	}
	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

// RunExec runs when the Exec command is invoked
func RunExec(out io.Writer, cmd *cobra.Command) error {
	return nil
}
