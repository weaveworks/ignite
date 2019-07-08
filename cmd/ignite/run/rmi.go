package run

import (
	"fmt"
	"os"

	"github.com/weaveworks/ignite/pkg/client"
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
		if image, err := client.Images().Find(filter.NewIDNameFilter(match)); err == nil {
			ro.images = append(ro.images, &imgmd.Image{image})
		} else {
			return nil, err
		}
	}

	if allVMs, err := client.VMs().FindAll(filter.NewAllFilter()); err == nil {
		ro.allVMs = make([]*vmmd.VM, 0, len(allVMs))
		for _, vm := range allVMs {
			ro.allVMs = append(ro.allVMs, &vmmd.VM{vm})
		}
	} else {
		return nil, err
	}

	return ro, nil
}

func Rmi(ro *rmiOptions) error {
	for _, image := range ro.images {
		for _, vm := range ro.allVMs {
			// Check if there's any VM using this image
			if vm.Spec.Image.UID == image.GetUID() {
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
