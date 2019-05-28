package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/metadata/kernmd"
	"github.com/luxas/ignite/pkg/util"
	"github.com/spf13/cobra"
	"io"
)

type kernelsOptions struct {
	kernels []*kernmd.KernelMetadata
}

// NewCmdKernels lists the available kernels
func NewCmdKernels(out io.Writer) *cobra.Command {
	ko := &kernelsOptions{}

	cmd := &cobra.Command{
		Use:   "kernels",
		Short: "List available kernels",
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if ko.kernels, err = matchAllKernels(); err != nil {
					return err
				}
				return RunKernels(ko)
			}())
		},
	}

	return cmd
}

func RunKernels(ko *kernelsOptions) error {
	o := util.NewOutput()
	defer o.Flush()

	o.Write("KERNEL ID", "CREATED", "SIZE", "NAME")
	for _, md := range ko.kernels {
		size, err := md.Size()
		if err != nil {
			return fmt.Errorf("failed to get size for %s %q: %v", md.Type, md.ID, err)
		}

		o.Write(md.ID, md.Created, util.ByteCountDecimal(size), md.Name)
	}

	return nil
}
