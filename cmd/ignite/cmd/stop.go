package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/vmcmd"
)

// NewCmdStop is an alias for vmcmd.NewCmdStop
func NewCmdStop(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdStop(out)
}
