package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/vmcmd"
)

// NewCmdRun is an alias for vmcmd.NewCmdRun
func NewCmdRun(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdRun(out)
}
