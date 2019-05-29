package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/cmd/vmcmd"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdRm is an alias for vmcmd.NewCmdRm
func NewCmdRm(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdRm(out)
}
