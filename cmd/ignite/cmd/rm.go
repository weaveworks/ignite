package cmd

import (
	"github.com/spf13/cobra"
	"io"
)

// NewCmdRm is an alias for NewCmdVMRm
func NewCmdRm(out io.Writer) *cobra.Command {
	return NewCmdVMRm(out)
}
