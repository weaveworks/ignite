package cmd

import (
	"github.com/spf13/cobra"
	"io"
)

// NewCmdRmi is an alias for NewCmdImageRm
func NewCmdRmi(out io.Writer) *cobra.Command {
	cmd := NewCmdImageRm(out)
	cmd.Use = "rmi [image]"

	return cmd
}
