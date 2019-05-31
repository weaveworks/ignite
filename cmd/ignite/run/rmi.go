package run

import (
	"fmt"
	"github.com/luxas/ignite/pkg/metadata/imgmd"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"os"
)

type RmiOptions struct {
	Image *imgmd.ImageMetadata
	VMs   []*vmmd.VMMetadata
}

func Rmi(ro *RmiOptions) error {
	for _, vm := range ro.VMs {
		if vm.VMOD().ImageID == ro.Image.ID {
			return fmt.Errorf("unable to remove, image %q is in use by VM %q", ro.Image.ID, vm.ID)
		}
	}

	if err := os.RemoveAll(ro.Image.ObjectPath()); err != nil {
		return fmt.Errorf("unable to remove directory for %s %q: %v", ro.Image.Type, ro.Image.ID, err)
	}

	fmt.Println(ro.Image.ID)
	return nil
}
