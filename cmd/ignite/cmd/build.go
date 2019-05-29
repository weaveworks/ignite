package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/cmd/imgcmd"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdBuild is an alias for imgcmd.NewCmdBuild
func NewCmdBuild(out io.Writer) *cobra.Command {
	return imgcmd.NewCmdBuild(out)
}
