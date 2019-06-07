package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/kerncmd"
)

// NewCmdRmk is an alias for kerncmd.NewCmdRm
func NewCmdRmk(out io.Writer) *cobra.Command {
	cmd := kerncmd.NewCmdRm(out)
	cmd.Use = "rmk <kernel>"

	return cmd
}
