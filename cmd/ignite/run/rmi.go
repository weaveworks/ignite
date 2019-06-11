package run

import (
	"fmt"
	"os"

	"github.com/weaveworks/ignite/cmd/ignite/run/runutil"

	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
)

type RmiFlags struct {
	Force bool
}

type rmiOptions struct {
	*RmiFlags
	images []*imgmd.ImageMetadata
	allVMs []*vmmd.VMMetadata
}

func (rf *RmiFlags) NewRmiOptions(l *runutil.ResLoader, imageMatches []string) (*rmiOptions, error) {
	ro := &rmiOptions{RmiFlags: rf}

	if allImages, err := l.Images(); err == nil {
		if ro.images, err = allImages.MatchMultiple(imageMatches); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	if allVMs, err := l.VMs(); err == nil {
		ro.allVMs = allVMs.MatchAll()
	} else {
		return nil, err
	}

	return ro, nil
}

func Rmi(ro *rmiOptions) error {
	for _, image := range ro.images {
		for _, vm := range ro.allVMs {
			// Check if there's any VM using this image
			if vm.VMOD().ImageID.Equal(image.ID) {
				if ro.Force {
					// Force-kill and remove the VM used by this image
					if err := Rm(&rmOptions{
						&RmFlags{Force: true},
						[]*vmmd.VMMetadata{vm},
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
