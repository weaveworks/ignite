package run

import (
	"fmt"

	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/filter"
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

func (rf *RmFlags) NewRmOptions(vmMatches []string) (*rmOptions, error) {
	ro := &rmOptions{RmFlags: rf}

	for _, match := range vmMatches {
		if vm, err := client.VMs().Find(filter.NewVMFilter(match)); err == nil {
			ro.vms = append(ro.vms, &vmmd.VM{vm})
		} else {
			return nil, err
		}
	}

	return ro, nil
}

func Rm(ro *rmOptions) error {
	for _, vm := range ro.vms {
		// If the VM is running, but we haven't enabled force-mode, return an error
		if vm.Running() && !ro.Force {
			return fmt.Errorf("%s is running", vm.GetUID())
		}

		// This will first kill the VM container, and then remove it
		if err := operations.RemoveVM(client.DefaultClient, vm.VM); err != nil {
			return err
		}
	}

	return nil
}
