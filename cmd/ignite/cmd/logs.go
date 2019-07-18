package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/vmcmd"
)

// NewCmdLogs is an alias for vmcmd.NewCmdLogs
func NewCmdLogs(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdLogs(out)
}
