package run

import (
	"fmt"
	"github.com/luxas/ignite/pkg/metadata/kernmd"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"os"
)

type RmkOptions struct {
	Kernel *kernmd.KernelMetadata
	VMs    []*vmmd.VMMetadata
}

func Rmk(ro *RmkOptions) error {
	for _, vm := range ro.VMs {
		if vm.VMOD().KernelID == ro.Kernel.ID {
			return fmt.Errorf("unable to remove, kernel %q is in use by VM %q", ro.Kernel.ID, vm.ID)
		}
	}

	if err := os.RemoveAll(ro.Kernel.ObjectPath()); err != nil {
		return fmt.Errorf("unable to remove directory for %s %q: %v", ro.Kernel.Type, ro.Kernel.ID, err)
	}

	fmt.Println(ro.Kernel.ID)
	return nil
}
