package run

import (
	"fmt"
	"github.com/luxas/ignite/pkg/metadata/kernmd"
	"os"
)

type RmkOptions struct {
	Kernel *kernmd.KernelMetadata
}

func Rmk(ro *RmkOptions) error {
	// TODO: Check that the given kernel is not used by any VMs
	if err := os.RemoveAll(ro.Kernel.ObjectPath()); err != nil {
		return fmt.Errorf("unable to remove directory for %s %q: %v", ro.Kernel.Type, ro.Kernel.ID, err)
	}

	fmt.Println(ro.Kernel.ID)
	return nil
}
