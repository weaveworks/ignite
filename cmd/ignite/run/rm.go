package run

import (
	"fmt"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	"github.com/weaveworks/ignite/pkg/config"
	"github.com/weaveworks/ignite/pkg/operations"
	"github.com/weaveworks/ignite/pkg/providers"
)

// RmFlags contains the flags supported by the remove command.
type RmFlags struct {
	Force      bool
	ConfigFile string
}

type RmOptions struct {
	*RmFlags
	vms []*api.VM
}

// NewRmOptions creates and returns RmOptions with all the flags and VMs to be
// removed.
func (rf *RmFlags) NewRmOptions(vmMatches []string) (*RmOptions, error) {
	ro := &RmOptions{RmFlags: rf}

	// If config file is provided, use it to find the VM to be removed.
	if len(rf.ConfigFile) != 0 {
		if len(vmMatches) > 0 {
			return ro, fmt.Errorf("cannot use both config flag and vm argument")
		}

		vm := &api.VM{}
		if err := scheme.Serializer.DecodeFileInto(rf.ConfigFile, vm); err != nil {
			return ro, err
		}
		// Name or UID must be provided in the config file.
		if len(vm.Name) == 0 && len(vm.UID) == 0 {
			return ro, fmt.Errorf("API resource config must have Name or UID")
		}
		ro.vms = []*api.VM{vm}
		return ro, nil
	}

	// Use vm args to find the VMs to be removed.
	if len(vmMatches) < 1 {
		return ro, fmt.Errorf("need at least one vm identifier as argument")
	}
	var err error
	ro.vms, err = getVMsForMatches(vmMatches)
	return ro, err
}

// Rm removes VMs based on RmOptions.
func Rm(ro *RmOptions) error {
	for _, vm := range ro.vms {
		// If the VM is running, but we haven't enabled force-mode, return an error
		if vm.Running() && !ro.Force {
			return fmt.Errorf("%s is running", vm.GetUID())
		}

		// Runtime and network info are present only when the VM is running.
		if vm.Running() {
			// Set the runtime and network-plugin providers from the VM status.
			if err := config.SetAndPopulateProviders(vm.Status.Runtime.Name, vm.Status.Network.Plugin); err != nil {
				return err
			}
		}

		// This will first kill the VM container, and then remove it
		if err := operations.DeleteVM(providers.Client, vm); err != nil {
			return err
		}
	}

	return nil
}
