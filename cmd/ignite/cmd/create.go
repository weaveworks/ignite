package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/vmcmd"
)

// NewCmdCreate is an alias for vmcmd.NewCmdCreate
func NewCmdCreate(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdCreate(out)
}
