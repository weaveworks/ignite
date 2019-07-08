package kerncmd

import (
	"io"

	"github.com/lithammer/dedent"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
)

// NewCmdKernel handles kernel-related functionality via its subcommands
// This command by itself lists available kernels
func NewCmdKernel(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kernel",
		Short: "Manage VM kernels",
		Long: dedent.Dedent(`
			Groups together functionality for managing VM kernels.
			Calling this command alone lists all available kernels.
		`),
		Aliases: []string{"kernels"},
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				ko, err := run.NewKernelsOptions()
				if err != nil {
					return err
				}

				return run.Kernels(ko)
			}())
		},
	}

	cmd.AddCommand(NewCmdLs(out))
	cmd.AddCommand(NewCmdRm(out))
	return cmd
}
