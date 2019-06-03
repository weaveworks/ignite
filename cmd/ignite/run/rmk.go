package run

import (
	"fmt"
	"github.com/luxas/ignite/pkg/metadata/kernmd"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"os"
)

type RmkOptions struct {
	Kernels []*kernmd.KernelMetadata
	VMs     []*vmmd.VMMetadata
}

func Rmk(ro *RmkOptions) error {
	for _, kernel := range ro.Kernels {
		for _, vm := range ro.VMs {
			if vm.VMOD().KernelID == kernel.ID {
				return fmt.Errorf("unable to remove, kernel %q is in use by VM %q", kernel.ID, vm.ID)
			}
		}

		if err := os.RemoveAll(kernel.ObjectPath()); err != nil {
			return fmt.Errorf("unable to remove directory for %s %q: %v", kernel.Type, kernel.ID, err)
		}

		fmt.Println(kernel.ID)
	}

	return nil
}
