package operations

import (
	"fmt"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"log"

	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/logs"
	"github.com/weaveworks/ignite/pkg/util"
)

var (
	stopArgs = []string{"stop"}
	killArgs = []string{"kill", "-s", "SIGQUIT"}
	rmArgs   = []string{"rm", "-f"}
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
	// Specify what container to remove
	dockerArgs := append(rmArgs, constants.IGNITE_PREFIX+vm.GetUID().String())

	// Remove the VM container
	// TODO: Use pkg/runtime here instead
	if _, err := util.ExecuteCommand("docker", dockerArgs...); err != nil {
		return fmt.Errorf("failed to remove container for VM %q: %v", vm.GetUID(), err)
	}

	return nil
}

// StopVM stops or kills a VM
func StopVM(vm *vmmd.VM, kill, silent bool) error {
	dockerArgs := stopArgs

	// Change to kill arguments if requested
	if kill {
		dockerArgs = killArgs
	}

	// Specify what container to stop/kill
	dockerArgs = append(dockerArgs, constants.IGNITE_PREFIX+vm.GetUID().String())

	// Stop/Kill the VM in docker
	// TODO: Use pkg/runtime here instead
	if _, err := util.ExecuteCommand("docker", dockerArgs...); err != nil {
		return fmt.Errorf("failed to stop container for VM %q: %v", vm.GetUID(), err)
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
