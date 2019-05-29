package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/cmd/vmcmd"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdStart is an alias for vmcmd.NewCmdStart
func NewCmdStart(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdStart(out)
}
