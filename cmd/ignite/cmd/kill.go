package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/vmcmd"
)

// NewCmdKill is an alias for vmcmd.NewCmdKill
func NewCmdKill(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdKill(out)
}
