package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

// NewCmdLogs is an alias for NewCmdVMLogs
func NewCmdLogs(out io.Writer) *cobra.Command {
	return NewCmdVMLogs(out)
}
