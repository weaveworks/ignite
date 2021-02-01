package run

import (
	"fmt"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/config"
	"github.com/weaveworks/ignite/pkg/providers"
)

// checkRunning can be used to skip the running check, this is used by Start and Run
// as the in-container ignite takes some time to start up and update the state
type AttachOptions struct {
	checkRunning bool
	vm           *api.VM
}

func NewAttachOptions(vmMatch string) (ao *AttachOptions, err error) {
	ao = &AttachOptions{checkRunning: true}
	ao.vm, err = getVMForMatch(vmMatch)
	return
}

func Attach(ao *AttachOptions) error {
	// Check if the VM is running
	if ao.checkRunning && !ao.vm.Running() {
		return fmt.Errorf("VM %q is not running", ao.vm.GetUID())
	}

	// Set the runtime and network-plugin providers from the VM status.
	if err := config.SetAndPopulateProviders(ao.vm.Status.Runtime.Name, ao.vm.Status.Network.Plugin); err != nil {
		return err
	}

	// Print the ID before attaching
	fmt.Println(ao.vm.GetUID())

	// Attach to the VM in Docker
	if err := providers.Runtime.AttachContainer(ao.vm.PrefixedID()); err != nil {
		return fmt.Errorf("failed to attach to container for VM %s: %v", ao.vm.GetUID(), err)
	}

	return nil
}
