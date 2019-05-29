package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/cmd/vmcmd"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdCreate is an alias for vmcmd.NewCmdCreate
func NewCmdCreate(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdCreate(out)
}
