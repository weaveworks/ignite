package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/cmd/vmcmd"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdStop is an alias for vmcmd.NewCmdStop
func NewCmdStop(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdStop(out)
}
