package run

import (
	"fmt"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"os"
)

type RmOptions struct {
	VM    *vmmd.VMMetadata
	Force bool
}

func Rm(ro *RmOptions) error {
	// Check if the VM is running
	if ro.VM.Running() {
		// If force is set, kill the VM
		if ro.Force {
			if err := Kill(&KillOptions{
				VM: ro.VM,
			}); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("%s is running", ro.VM.ID)
		}
	}

	if err := os.RemoveAll(ro.VM.ObjectPath()); err != nil {
		return fmt.Errorf("unable to remove directory for %s %q: %v", ro.VM.Type, ro.VM.ID, err)
	}

	fmt.Println(ro.VM.ID)
	return nil
}
