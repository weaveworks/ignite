package run

import (
	"fmt"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	"github.com/weaveworks/ignite/pkg/operations"
	"github.com/weaveworks/ignite/pkg/providers"
)

// RmFlags contains the flags supported by the remove command.
type RmFlags struct {
	Force      bool
	ConfigFile string
}

type rmOptions struct {
	*RmFlags
	vms []*api.VM
}

// NewRmOptions creates and returns rmOptions with all the flags and VMs to be
// removed.
func (rf *RmFlags) NewRmOptions(vmMatches []string) (*rmOptions, error) {
	ro := &rmOptions{RmFlags: rf}

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

// Rm removes VMs based on rmOptions.
func Rm(ro *rmOptions) error {
	for _, vm := range ro.vms {
		// If the VM is running, but we haven't enabled force-mode, return an error
		if vm.Running() && !ro.Force {
			return fmt.Errorf("%s is running", vm.GetUID())
		}

		// This will first kill the VM container, and then remove it
		if err := operations.DeleteVM(providers.Client, vm); err != nil {
			return err
		}
	}

	return nil
}
