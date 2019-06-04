package run

import (
	"fmt"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"os"
)

type RmOptions struct {
	VMs   []*vmmd.VMMetadata
	Force bool
}

func Rm(ro *RmOptions) error {
	for _, vm := range ro.VMs {
		// Check if the VM is running
		if vm.Running() {
			// If force is set, kill the VM
			if ro.Force {
				if err := Stop(&StopOptions{
					VMs:  []*vmmd.VMMetadata{vm},
					Kill: true,
				}); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("%v is running", vm)
			}
		}

		if err := os.RemoveAll(vm.ObjectPath()); err != nil {
			return fmt.Errorf("unable to remove directory for %s %q: %v", vm.Type, vm.ID, err)
		}

		fmt.Println(vm.ID)
	}

	return nil
}
