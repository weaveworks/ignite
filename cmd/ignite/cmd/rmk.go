package cmd

import (
	"github.com/spf13/cobra"
	"io"
)

// NewCmdRmk is an alias for NewCmdKernelRm
func NewCmdRmk(out io.Writer) *cobra.Command {
	cmd := NewCmdKernelRm(out)
	cmd.Use = "rmk [kernel]"

	return cmd
}
