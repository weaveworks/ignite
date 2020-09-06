package run

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/operations/lookup"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/libgitops/pkg/filter"
)

type RmkFlags struct {
	Force bool
}

type RmkOptions struct {
	*RmkFlags
	kernels []*api.Kernel
	allVMs  []*api.VM
}

func (rf *RmkFlags) NewRmkOptions(kernelMatches []string) (*RmkOptions, error) {
	ro := &RmkOptions{RmkFlags: rf}

	for _, match := range kernelMatches {
		if kernel, err := providers.Client.Kernels().Find(filter.NewIDNameFilter(match)); err == nil {
			ro.kernels = append(ro.kernels, kernel)
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

func Rmk(ro *RmkOptions) error {
	for _, kernel := range ro.kernels {
		for _, vm := range ro.allVMs {
			kernelUID, err := lookup.KernelUIDForVM(vm, providers.Client)
			if err != nil {
				log.Warnf("Could not lookup kernel UID for VM %q: %v", vm.GetUID(), err)
				continue
			}

			// Check if there's any VM using this kernel
			if kernelUID == kernel.GetUID() {
				if ro.Force {
					// Force-kill and remove the VM used by this kernel
					if err := Rm(&RmOptions{
						&RmFlags{Force: true},
						[]*api.VM{vm},
					}); err != nil {
						return err
					}
				} else {
					return fmt.Errorf("unable to remove, kernel %q is in use by VM %q", kernel.GetUID(), vm.GetUID())
				}
			}
		}

		if err := os.RemoveAll(kernel.ObjectPath()); err != nil {
			return fmt.Errorf("unable to remove directory for %s %q: %v", kernel.GetKind(), kernel.GetUID(), err)
		}

		fmt.Println(kernel.GetUID())
	}

	return nil
}
