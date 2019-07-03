package run

import (
	"fmt"
	"os"

	"github.com/weaveworks/ignite/pkg/metadata/loader"

	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
)

type RmkFlags struct {
	Force bool
}

type rmkOptions struct {
	*RmkFlags
	kernels []*kernmd.Kernel
	allVMs  []*vmmd.VMMetadata
}

func (rf *RmkFlags) NewRmkOptions(l *loader.ResLoader, kernelMatches []string) (*rmkOptions, error) {
	ro := &rmkOptions{RmkFlags: rf}

	if allKernels, err := l.Kernels(); err == nil {
		if ro.kernels, err = allKernels.MatchMultiple(kernelMatches); err != nil {
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

func Rmk(ro *rmkOptions) error {
	for _, kernel := range ro.kernels {
		for _, vm := range ro.allVMs {
			// Check if there's any VM using this kernel
			if vm.Spec.Kernel.ID == kernel.GetUID() {
				if ro.Force {
					// Force-kill and remove the VM used by this kernel
					if err := Rm(&rmOptions{
						&RmFlags{Force: true},
						[]*vmmd.VMMetadata{vm},
					}); err != nil {
						return err
					}
				} else {
					return fmt.Errorf("unable to remove, kernel %q is in use by VM %q", kernel.GetUID(), vm.GetUID())
				}
			}
		}

		if err := os.RemoveAll(kernel.ObjectPath()); err != nil {
			return fmt.Errorf("unable to remove directory for %s %q: %v", kernel.Type(), kernel.GetUID(), err)
		}

		fmt.Println(kernel.GetUID())
	}

	return nil
}
