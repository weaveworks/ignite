package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/vmcmd"
)

// NewCmdPreflight is an alias for vmcmd.NewCmdPreflight
func NewCmdPreflight(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdPreflight(out)
}
