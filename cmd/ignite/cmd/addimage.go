package cmd

import (
	"github.com/spf13/cobra"
	"io"
)

// NewCmdAddImage is an alias for NewCmdImageImport
func NewCmdAddImage(out io.Writer) *cobra.Command {
	cmd := NewCmdImageImport(out)
	cmd.Use = "addimage [path]"

	return cmd
}
