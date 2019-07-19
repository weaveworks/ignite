package run

import (
	"fmt"
	"os"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/filter"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
)

type RmkFlags struct {
	Force bool
}

type rmkOptions struct {
	*RmkFlags
	kernels []*api.Kernel
	allVMs  []*vmmd.VM
}

func (rf *RmkFlags) NewRmkOptions(kernelMatches []string) (*rmkOptions, error) {
	ro := &rmkOptions{RmkFlags: rf}

	for _, match := range kernelMatches {
		if kernel, err := client.Kernels().Find(filter.NewIDNameFilter(match)); err == nil {
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

func Rmk(ro *rmkOptions) error {
	for _, kernel := range ro.kernels {
		for _, vm := range ro.allVMs {
			// Check if there's any VM using this kernel
			if vm.GetKernelUID() == kernel.GetUID() {
				if ro.Force {
					// Force-kill and remove the VM used by this kernel
					if err := Rm(&rmOptions{
						&RmFlags{Force: true},
						[]*vmmd.VM{vm},
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
