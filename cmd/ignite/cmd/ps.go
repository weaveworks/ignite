package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/cmd/vmcmd"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdPs is an alias for vmcmd.NewCmdPs
func NewCmdPs(out io.Writer) *cobra.Command {
	cmd := vmcmd.NewCmdPs(out)
	cmd.Aliases = nil
	return cmd
}
