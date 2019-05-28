package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/metadata"
	"github.com/luxas/ignite/pkg/metadata/imgmd"
	"github.com/luxas/ignite/pkg/metadata/kernmd"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"github.com/luxas/ignite/pkg/util"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdPs lists running Firecracker VMs
func NewCmdPs(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ps",
		Short: "List running Firecracker VMs",
		Run: func(cmd *cobra.Command, args []string) {
			err := RunPs(out, cmd)
			errutils.Check(err)
		},
	}

	cmd.Flags().BoolVarP(&all, "all", "a", false, "Show all VMs instead of just running ones")
	return cmd
}

func RunPs(out io.Writer, cmd *cobra.Command) error {
	var mds []*vmmd.VMMetadata

	// Match all VMs using the VMFilter with state checking
	if matches, err := filter.NewFilterer(vmmd.NewVMFilterAll("", all), metadata.VM.Path(), vmmd.LoadVMMetadata); err == nil {
		if all, err := matches.All(); err == nil {
			if mds, err = vmmd.ToVMMetadataAll(all); err != nil {
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

	o.Write("VM ID", "IMAGE", "KERNEL", "CREATED", "SIZE", "CPUS", "MEMORY", "STATE", "NAME")
	for _, md := range mds {
		od := md.ObjectData.(*vmmd.VMObjectData)
		size, err := md.Size()
		if err != nil {
			return fmt.Errorf("failed to get size for %s %q: %v", md.Type, md.ID, err)
		}

		image := imgmd.NewImageMetadata(od.ImageID, "")
		kernel := kernmd.NewKernelMetadata(od.KernelID, "")

		if err := image.Load(); err != nil {
			return fmt.Errorf("failed to load image metadata for %s %q: %v", md.Type, md.ID, err)
		}

		if err := kernel.Load(); err != nil {
			return fmt.Errorf("failed to load kernel metadata for %s %q: %v", md.Type, md.ID, err)
		}

		o.Write(md.ID, image.Name, kernel.Name, md.Created, util.ByteCountDecimal(size), od.VCPUs, util.ByteCountDecimal(od.Memory*1000000), od.State, md.Name)
	}

	return nil
}
