package cmd

import (
	"github.com/spf13/cobra"
	"io"
)

// NewCmdStart is an alias for NewCmdVMStart
func NewCmdStart(out io.Writer) *cobra.Command {
	return NewCmdVMStart(out)
}
