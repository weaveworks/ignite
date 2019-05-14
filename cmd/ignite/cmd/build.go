package cmd

import (
	"io"

	"github.com/luxas/ignite/pkg/errors"
	"github.com/spf13/cobra"
)

// NewCmdBuild builds a Firecracker VM.
func NewCmdBuild(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build a Firecracker VM",
		Run: func(cmd *cobra.Command, args []string) {
			err := RunBuild(out, cmd)
			errors.Check(err)
		},
	}
	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

// RunBuild runs when the Build command is invoked
func RunBuild(out io.Writer, cmd *cobra.Command) error {
	return nil
}
