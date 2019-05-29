package cmd

import (
	"github.com/spf13/cobra"
	"io"
)

// NewCmdPs is an alias for NewCmdVMPs
func NewCmdPs(out io.Writer) *cobra.Command {
	cmd := NewCmdVMPs(out)
	cmd.Aliases = nil
	return cmd
}
