package run

import (
	"fmt"
	"os"

	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
)

type RmkOptions struct {
	Kernels []*kernmd.KernelMetadata
	VMs     []*vmmd.VMMetadata
	Force   bool
}

func Rmk(ro *RmkOptions) error {
	for _, kernel := range ro.Kernels {
		for _, vm := range ro.VMs {
			// Check if there's any VM using this kernel
			if vm.VMOD().KernelID == kernel.ID {
				if ro.Force {
					// Force-kill and remove the VM used by this kernel
					if err := Rm(&RmOptions{
						VMs:   []*vmmd.VMMetadata{vm},
						Force: true,
					}); err != nil {
						return err
					}
				} else {
					return fmt.Errorf("unable to remove, kernel %q is in use by VM %q", kernel.ID, vm.ID)
				}
			}
		}

		if err := os.RemoveAll(kernel.ObjectPath()); err != nil {
			return fmt.Errorf("unable to remove directory for %s %q: %v", kernel.Type, kernel.ID, err)
		}

		fmt.Println(kernel.ID)
	}

	return nil
}
