package cmd

import (
	"github.com/spf13/cobra"
	"io"
)

// NewCmdBuild is an alias for NewCmdImageBuild
func NewCmdBuild(out io.Writer) *cobra.Command {
	return NewCmdImageBuild(out)
}
