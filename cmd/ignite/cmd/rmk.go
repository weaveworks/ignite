package cmd

import (
	"io"

	"github.com/weaveworks/ignite/cmd/ignite/cmd/kerncmd"
	"github.com/spf13/cobra"
)

// NewCmdRmk is an alias for kerncmd.NewCmdRm
func NewCmdRmk(out io.Writer) *cobra.Command {
	cmd := kerncmd.NewCmdRm(out)
	cmd.Use = "rmk [kernel]"

	return cmd
}
