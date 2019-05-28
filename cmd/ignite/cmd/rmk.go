package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/metadata/kernmd"
	"github.com/spf13/cobra"
	"io"
	"os"
)

type rmkOptions struct {
	kernel *kernmd.KernelMetadata
}

// NewCmdRmk removes the given kernel
func NewCmdRmk(out io.Writer) *cobra.Command {
	ro := &rmkOptions{}

	cmd := &cobra.Command{
		Use:   "rmk [id]",
		Short: "Remove a kernel",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if ro.kernel, err = matchSingleKernel(args[0]); err != nil {
					return err
				}
				return RunRmk(ro)
			}())
		},
	}
	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

func RunRmk(ro *rmkOptions) error {
	// TODO: Check that the given kernel is not used by any VMs
	if err := os.RemoveAll(ro.kernel.ObjectPath()); err != nil {
		return fmt.Errorf("unable to remove directory for %s %q: %v", ro.kernel.Type, ro.kernel.ID, err)
	}

	fmt.Println(ro.kernel.ID)
	return nil
}
