package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/cmd/kerncmd"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdAddKernel is an alias for kerncmd.NewCmdImport
func NewCmdAddKernel(out io.Writer) *cobra.Command {
	cmd := kerncmd.NewCmdImport(out)
	cmd.Use = "addkernel [path]"

	return cmd
}
