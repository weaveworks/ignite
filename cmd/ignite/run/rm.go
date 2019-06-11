package run

import (
	"fmt"

	"github.com/weaveworks/ignite/cmd/ignite/run/runutil"
	"github.com/weaveworks/ignite/pkg/logs"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
)

type RmFlags struct {
	Force bool
}

type rmOptions struct {
	*RmFlags
	vms []*vmmd.VMMetadata
}

func (rf *RmFlags) NewRmOptions(l *runutil.ResLoader, vmMatches []string) (*rmOptions, error) {
	ro := &rmOptions{RmFlags: rf}

	if allVMs, err := l.VMs(); err == nil {
		if ro.vms, err = allVMs.MatchMultiple(vmMatches); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return ro, nil
}

func Rm(ro *rmOptions) error {
	for _, vm := range ro.vms {
		// Check if the VM is running
		if vm.Running() {
			// If force is set, kill the vm
			if ro.Force {
				if err := Stop(&stopOptions{
					&StopFlags{
						Kill: true,
					},
					[]*vmmd.VMMetadata{vm},
					true,
				}); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("%v is running", vm)
			}
		}

		if err := vm.Remove(logs.Quiet); err != nil {
			return err
		}
	}

	return nil
}
