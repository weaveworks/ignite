package run

import (
	"fmt"

	"github.com/weaveworks/ignite/pkg/logs"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
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
					VMs:          []*vmmd.VMMetadata{vm},
					Kill:         true,
					DisablePrint: true,
				}); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("%v is running", vm)
			}
		}

		if err := vm.Remove(logs.Quiet); err != nil {
			return err
		}
	}

	return nil
}
