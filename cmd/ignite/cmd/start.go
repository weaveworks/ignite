package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/vmcmd"
)

// NewCmdStart is an alias for vmcmd.NewCmdStart
func NewCmdStart(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdStart(out)
}
