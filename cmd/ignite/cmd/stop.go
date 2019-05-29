package cmd

import (
	"github.com/spf13/cobra"
	"io"
)

// NewCmdStop is an alias for NewCmdVMStop
func NewCmdStop(out io.Writer) *cobra.Command {
	return NewCmdVMStop(out)
}
