package run

import (
	"fmt"
	"github.com/weaveworks/ignite/pkg/providers"
	"os"

	"github.com/weaveworks/ignite/pkg/filter"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
)

type RmiFlags struct {
	Force bool
}

type rmiOptions struct {
	*RmiFlags
	images []*imgmd.Image
	allVMs []*vmmd.VM
}

func (rf *RmiFlags) NewRmiOptions(imageMatches []string) (*rmiOptions, error) {
	ro := &rmiOptions{RmiFlags: rf}

	for _, match := range imageMatches {
		if image, err := providers.Client.Images().Find(filter.NewIDNameFilter(match)); err == nil {
			ro.images = append(ro.images, imgmd.WrapImage(image))
		} else {
			return nil, err
		}
	}

	var err error
	ro.allVMs, err = getAllVMs()
	if err != nil {
		return nil, err
	}

	return ro, nil
}

func Rmi(ro *rmiOptions) error {
	for _, image := range ro.images {
		for _, vm := range ro.allVMs {
			// Check if there's any VM using this image
			if vm.GetImageUID() == image.GetUID() {
				if ro.Force {
					// Force-kill and remove the VM used by this image
					if err := Rm(&rmOptions{
						&RmFlags{Force: true},
						[]*vmmd.VM{vm},
					}); err != nil {
						return err
					}
				} else {
					return fmt.Errorf("unable to remove, image %q is in use by VM %q", image.GetUID(), vm.GetUID())
				}
			}
		}

		if err := os.RemoveAll(image.ObjectPath()); err != nil {
			return fmt.Errorf("unable to remove directory for %s %q: %v", image.GetKind(), image.GetUID(), err)
		}

		fmt.Println(image.GetUID())
	}

	return nil
}
