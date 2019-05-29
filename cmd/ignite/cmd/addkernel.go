package cmd

import (
	"github.com/spf13/cobra"
	"io"
)

// NewCmdAddKernel is an alias for NewCmdKernelImport
func NewCmdAddKernel(out io.Writer) *cobra.Command {
	cmd := NewCmdKernelImport(out)
	cmd.Use = "addkernel [path]"

	return cmd
}
