package cmd

import (
	"github.com/spf13/cobra"
	"io"
)

// NewCmdAttach is an alias for NewCmdVMAttach
func NewCmdAttach(out io.Writer) *cobra.Command {
	return NewCmdVMAttach(out)
}
