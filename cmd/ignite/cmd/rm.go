package cmd

import (
	"io"

	"github.com/weaveworks/ignite/cmd/ignite/cmd/vmcmd"
	"github.com/spf13/cobra"
)

// NewCmdRm is an alias for vmcmd.NewCmdRm
func NewCmdRm(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdRm(out)
}
