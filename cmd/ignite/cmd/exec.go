package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/vmcmd"
)

// NewCmdExec is an alias for vmcmd.NewCmdExec
func NewCmdExec(out io.Writer, err io.Writer, in io.Reader) *cobra.Command {
	return vmcmd.NewCmdExec(out, err, in)
}
