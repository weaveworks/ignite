package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/vmcmd"
)

// NewCmdSSH is an alias for vmcmd.NewCmdSSH
func NewCmdSSH(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdSSH(out)
}
