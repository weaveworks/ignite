package run

import (
	"fmt"
	"github.com/luxas/ignite/pkg/metadata/imgmd"
	"os"
)

type RmiOptions struct {
	Image *imgmd.ImageMetadata
}

func Rmi(ro *RmiOptions) error {
	// TODO: Check that the given image is not used by any VMs
	if err := os.RemoveAll(ro.Image.ObjectPath()); err != nil {
		return fmt.Errorf("unable to remove directory for %s %q: %v", ro.Image.Type, ro.Image.ID, err)
	}

	fmt.Println(ro.Image.ID)
	return nil
}
