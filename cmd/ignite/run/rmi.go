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
	Force  bool
}

func Rmi(ro *RmiOptions) error {
	for _, image := range ro.Images {
		for _, vm := range ro.VMs {
			// Check if there's any VM using this image
			if vm.VMOD().ImageID == image.ID {
				if ro.Force {
					// Force-kill and remove the VM used by this image
					if err := Rm(&RmOptions{
						VMs:   []*vmmd.VMMetadata{vm},
						Force: true,
					}); err != nil {
						return err
					}
				} else {
					return fmt.Errorf("unable to remove, image %q is in use by VM %q", image.ID, vm.ID)
				}
			}
		}

		if err := os.RemoveAll(image.ObjectPath()); err != nil {
			return fmt.Errorf("unable to remove directory for %s %q: %v", image.Type, image.ID, err)
		}

		fmt.Println(image.ID)
	}

	return nil
}
