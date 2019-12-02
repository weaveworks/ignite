package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/vmcmd"
)

// NewCmdCP is an alias for vmcmd.NewCmdCP
func NewCmdCP(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdCP(out)
}
