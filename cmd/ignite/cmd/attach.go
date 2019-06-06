package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/vmcmd"
)

// NewCmdAttach is an alias for vmcmd.NewCmdAttach
func NewCmdAttach(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdAttach(out)
}
