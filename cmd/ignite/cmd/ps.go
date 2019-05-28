package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/metadata/imgmd"
	"github.com/luxas/ignite/pkg/metadata/kernmd"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"github.com/luxas/ignite/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
)

type psOptions struct {
	vms []*vmmd.VMMetadata
	all bool
}

// NewCmdPs lists running Firecracker VMs
func NewCmdPs(out io.Writer) *cobra.Command {
	po := &psOptions{}

	cmd := &cobra.Command{
		Use:   "ps",
		Short: "List running Firecracker VMs",
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if po.vms, err = matchAllVMs(po.all); err != nil {
					return err
				}
				return RunPs(po)
			}())
		},
	}

	addPsFlags(cmd.Flags(), po)
	return cmd
}

func addPsFlags(fs *pflag.FlagSet, po *psOptions) {
	fs.BoolVarP(&po.all, "all", "a", false, "Show all VMs, not just running ones")
}

func RunPs(po *psOptions) error {
	o := util.NewOutput()
	defer o.Flush()

	o.Write("VM ID", "IMAGE", "KERNEL", "CREATED", "SIZE", "CPUS", "MEMORY", "STATE", "NAME")
	for _, vm := range po.vms {
		od := vm.ObjectData.(*vmmd.VMObjectData)
		size, err := vm.Size()
		if err != nil {
			return fmt.Errorf("failed to get size for %s %q: %v", vm.Type, vm.ID, err)
		}

		image := imgmd.NewImageMetadata(od.ImageID, "")
		kernel := kernmd.NewKernelMetadata(od.KernelID, "")

		if err := image.Load(); err != nil {
			return fmt.Errorf("failed to load image metadata for %s %q: %v", vm.Type, vm.ID, err)
		}

		if err := kernel.Load(); err != nil {
			return fmt.Errorf("failed to load kernel metadata for %s %q: %v", vm.Type, vm.ID, err)
		}

		// TODO: Clean up this print
		o.Write(vm.ID, image.Name, kernel.Name, vm.Created, util.ByteCountDecimal(size), od.VCPUs, util.ByteCountDecimal(od.Memory*1000000), od.State, vm.Name)
	}

	return nil
}
