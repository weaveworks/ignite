package cmd

import (
	"github.com/spf13/cobra"
	"io"
)

// NewCmdKill is an alias for NewCmdVMKill
func NewCmdKill(out io.Writer) *cobra.Command {
	return NewCmdVMKill(out)
}
