package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/imgcmd"
)

// NewCmdBuild is an alias for imgcmd.NewCmdBuild
func NewCmdBuild(out io.Writer) *cobra.Command {
	return imgcmd.NewCmdBuild(out)
}
