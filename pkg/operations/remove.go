package operations

import (
	"fmt"
	"github.com/weaveworks/ignite/pkg/providers"
	"log"

	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/logs"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/util"
)

const (
	signalSIGQUIT = "SIGQUIT"
)

// RemoveVM removes the specified VM
func RemoveVM(c *client.Client, vm *vmmd.VM) error {
	// If the VM is running, try to kill it first so we don't leave dangling containers
	if vm.Running() {
		if err := StopVM(vm, true, true); err != nil {
			return err
		}
	}

	if err := c.VMs().Delete(vm.GetUID()); err != nil {
		return err
	}

	// Force-remove the VM container. Don't care about the error.
	_ = RemoveVMContainer(vm)

	if logs.Quiet {
		fmt.Println(vm.GetUID())
	} else {
		log.Printf("Removed %s with name %q and ID %q", vm.GetKind(), vm.GetName(), vm.GetUID())
	}

	return nil
}

func RemoveVMContainer(vm meta.Object) error {
	// Remove the VM container
	if err := providers.Runtime.RemoveContainer(util.NewPrefixer().Prefix(vm.GetUID())); err != nil {
		return fmt.Errorf("failed to remove container for VM %q: %v", vm.GetUID(), err)
	}

	return nil
}

// StopVM stops or kills a VM
func StopVM(vm *vmmd.VM, kill, silent bool) error {
	var err error
	container := util.NewPrefixer().Prefix(vm.GetUID())
	action := "stop"

	// Stop or kill the VM container
	if kill {
		action = "kill"
		err = providers.Runtime.KillContainer(container, signalSIGQUIT) // TODO: common constant for SIGQUIT
	} else {
		err = providers.Runtime.StopContainer(container, nil)
	}

	if err != nil {
		return fmt.Errorf("failed to %s container for %s %q: %v", action, vm.GetKind(), vm.GetUID(), err)
	}

	if silent {
		return nil
	}

	if logs.Quiet {
		fmt.Println(vm.GetUID())
	} else {
		log.Printf("Stopped %s with name %q and ID %q", vm.GetKind(), vm.GetName(), vm.GetUID())
	}

	return nil
}
