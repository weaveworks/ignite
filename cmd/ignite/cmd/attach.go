package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/cmd/vmcmd"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdAttach is an alias for vmcmd.NewCmdAttach
func NewCmdAttach(out io.Writer) *cobra.Command {
	return vmcmd.NewCmdAttach(out)
}
