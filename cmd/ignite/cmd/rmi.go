package cmd

import (
	"io"

	"github.com/luxas/ignite/cmd/ignite/cmd/imgcmd"
	"github.com/spf13/cobra"
)

// NewCmdRmi is an alias for imgcmd.NewCmdRm
func NewCmdRmi(out io.Writer) *cobra.Command {
	cmd := imgcmd.NewCmdRm(out)
	cmd.Use = "rmi [image]"

	return cmd
}
