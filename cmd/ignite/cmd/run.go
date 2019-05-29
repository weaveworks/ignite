package cmd

import (
	"github.com/spf13/cobra"
	"io"
)

// NewCmdRun is an alias for
func NewCmdRun(out io.Writer) *cobra.Command {
	return NewCmdVMRun(out)
}
