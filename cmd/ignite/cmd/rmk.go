package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/metadata"
	"github.com/luxas/ignite/pkg/metadata/kernmd"
	"github.com/spf13/cobra"
	"io"
	"os"
)

// NewCmdRmk removes the given kernel
func NewCmdRmk(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rmk [id]",
		Short: "Remove a kernel",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := RunRmk(out, cmd, args)
			errutils.Check(err)
		},
	}
	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

func RunRmk(out io.Writer, cmd *cobra.Command, args []string) error {
	var kernel *kernmd.KernelMetadata

	// Match a single Kernel using the KernelFilter
	if matches, err := filter.NewFilterer(kernmd.NewKernelFilter(args[1]), metadata.Kernel.Path(), kernmd.LoadKernelMetadata); err == nil {
		if filterable, err := matches.Single(); err == nil {
			if kernel, err = kernmd.ToKernelMetadata(filterable); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}

	// TODO: Check that the given kernel is not used by any VMs
	if err := os.RemoveAll(kernel.ObjectPath()); err != nil {
		return fmt.Errorf("unable to remove directory for %s %q: %v", kernel.Type, kernel.ID, err)
	}

	fmt.Println(kernel.ID)
	return nil
}
