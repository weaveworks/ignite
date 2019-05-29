package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/cmd/imgcmd"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdAddImage is an alias for imgcmd.NewCmdImport
func NewCmdAddImage(out io.Writer) *cobra.Command {
	cmd := imgcmd.NewCmdImport(out)
	cmd.Use = "addimage [path]"

	return cmd
}
