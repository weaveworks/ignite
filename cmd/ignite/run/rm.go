package run

import (
	"fmt"

	"github.com/weaveworks/ignite/pkg/providers"

	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/operations"
)

type RmFlags struct {
	Force bool
}

type rmOptions struct {
	*RmFlags
	vms []*vmmd.VM
}

func (rf *RmFlags) NewRmOptions(vmMatches []string) (ro *rmOptions, err error) {
	ro = &rmOptions{RmFlags: rf}
	ro.vms, err = getVMsForMatches(vmMatches)
	return
}

func Rm(ro *rmOptions) error {
	for _, vm := range ro.vms {
		// If the VM is running, but we haven't enabled force-mode, return an error
		if vm.Running() && !ro.Force {
			return fmt.Errorf("%s is running", vm.GetUID())
		}

		// This will first kill the VM container, and then remove it
		if err := operations.RemoveVM(providers.Client, vm); err != nil {
			return err
		}
	}

	return nil
}
