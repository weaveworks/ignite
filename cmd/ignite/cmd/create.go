package cmd

import (
	"github.com/spf13/cobra"
	"io"
)

// NewCmdCreate is an alias for NewCmdVMCreate
func NewCmdCreate(out io.Writer) *cobra.Command {
	return NewCmdVMCreate(out)
}
