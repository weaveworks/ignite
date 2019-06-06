package cmd

import (
	"io"

	"github.com/luxas/ignite/cmd/ignite/cmd/vmcmd"
	"github.com/spf13/cobra"
)

// NewCmdAttach is an alias for vmcmd.NewCmdAttach
func NewCmdAttach(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdAttach(out)
}
