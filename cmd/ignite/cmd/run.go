package cmd

import (
	"io"

	"github.com/weaveworks/ignite/cmd/ignite/cmd/vmcmd"
	"github.com/spf13/cobra"
)

// NewCmdRun is an alias for vmcmd.NewCmdRun
func NewCmdRun(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdRun(out)
}
