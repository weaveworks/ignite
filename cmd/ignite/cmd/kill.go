package cmd

import (
	"io"

	"github.com/luxas/ignite/cmd/ignite/cmd/vmcmd"
	"github.com/spf13/cobra"
)

// NewCmdKill is an alias for vmcmd.NewCmdKill
func NewCmdKill(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdKill(out)
}
