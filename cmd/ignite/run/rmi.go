package run

import (
	"fmt"
	"github.com/luxas/ignite/pkg/metadata/imgmd"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"os"
)

type RmiOptions struct {
	Images []*imgmd.ImageMetadata
	VMs    []*vmmd.VMMetadata
}

func Rmi(ro *RmiOptions) error {
	for _, image := range ro.Images {
		for _, vm := range ro.VMs {
			if vm.VMOD().ImageID == image.ID {
				return fmt.Errorf("unable to remove, image %q is in use by VM %q", image.ID, vm.ID)
			}
		}

		if err := os.RemoveAll(image.ObjectPath()); err != nil {
			return fmt.Errorf("unable to remove directory for %s %q: %v", image.Type, image.ID, err)
		}

		fmt.Println(image.ID)
	}

	return nil
}
