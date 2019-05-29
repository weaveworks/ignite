package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/cmd/vmcmd"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdRun is an alias for vmcmd.NewCmdRun
func NewCmdRun(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdRun(out)
}
