package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/metadata"
	"github.com/luxas/ignite/pkg/metadata/kernmd"
	"github.com/luxas/ignite/pkg/util"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdKernels lists the available kernels
func NewCmdKernels(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kernels",
		Short: "List available kernels",
		Run: func(cmd *cobra.Command, args []string) {
			err := RunKernels(out, cmd)
			errutils.Check(err)
		},
	}

	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

func RunKernels(out io.Writer, cmd *cobra.Command) error {
	var mds []*kernmd.KernelMetadata

	// Match all Kernels using the KernelFilter
	if matches, err := filter.NewFilterer(kernmd.NewKernelFilter(""), metadata.Kernel.Path(), kernmd.LoadKernelMetadata); err == nil {
		if all, err := matches.All(); err == nil {
			if mds, err = kernmd.ToKernelMetadataAll(all); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}

	o := util.NewOutput()
	defer o.Flush()

	o.Write("KERNEL ID", "CREATED", "SIZE", "NAME")
	for _, md := range mds {
		size, err := md.Size()
		if err != nil {
			return fmt.Errorf("failed to get size for %s %q: %v", md.Type, md.ID, err)
		}

		o.Write(md.ID, md.Created, size, md.Name)
	}

	return nil
}
