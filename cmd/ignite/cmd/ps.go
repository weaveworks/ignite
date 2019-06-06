package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/vmcmd"
)

// NewCmdPs is an alias for vmcmd.NewCmdPs
func NewCmdPs(out io.Writer) *cobra.Command {
	cmd := vmcmd.NewCmdPs(out)
	cmd.Aliases = nil
	return cmd
}
