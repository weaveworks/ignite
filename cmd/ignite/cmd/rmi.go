package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/imgcmd"
)

// NewCmdRmi is an alias for imgcmd.NewCmdRm
func NewCmdRmi(out io.Writer) *cobra.Command {
	cmd := imgcmd.NewCmdRm(out)
	cmd.Use = "rmi <image>"

	return cmd
}
