package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/cmd/vmcmd"
	"io"

	"github.com/spf13/cobra"
)

// NewCmdLogs is an alias for vmcmd.NewCmdLogs
func NewCmdLogs(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdLogs(out)
}
