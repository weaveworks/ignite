package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/vmcmd"
)

// NewCmdRm is an alias for vmcmd.NewCmdRm
func NewCmdRm(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdRm(out)
}
