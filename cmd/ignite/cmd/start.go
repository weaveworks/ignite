package cmd

import (
	"io"

	"github.com/luxas/ignite/pkg/errors"
	"github.com/spf13/cobra"
)

// NewCmdStart starts a Firecracker VM.
func NewCmdStart(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start a Firecracker VM",
		Run: func(cmd *cobra.Command, args []string) {
			err := RunStart(out, cmd)
			errors.Check(err)
		},
	}
	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

// RunStart runs when the start command is invoked
func RunStart(out io.Writer, cmd *cobra.Command) error {
	return nil
}
