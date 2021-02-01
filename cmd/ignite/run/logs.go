package run

import (
	"fmt"
	"io/ioutil"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/config"
	"github.com/weaveworks/ignite/pkg/providers"
)

type LogsOptions struct {
	vm *api.VM
}

func NewLogsOptions(vmMatch string) (lo *LogsOptions, err error) {
	lo = &LogsOptions{}
	lo.vm, err = getVMForMatch(vmMatch)
	return
}

func Logs(lo *LogsOptions) error {
	// Check if the VM is running
	if !lo.vm.Running() {
		return fmt.Errorf("VM %q is not running", lo.vm.GetUID())
	}

	// Set the runtime and network-plugin providers from the VM status.
	if err := config.SetAndPopulateProviders(lo.vm.Status.Runtime.Name, lo.vm.Status.Network.Plugin); err != nil {
		return err
	}

	// Fetch the VM logs
	rc, err := providers.Runtime.ContainerLogs(lo.vm.PrefixedID())
	if err != nil {
		return fmt.Errorf("failed to get logs for VM %q: %v", lo.vm.GetUID(), err)
	}

	// Read the stream to a byte buffer
	b, err := ioutil.ReadAll(rc)
	if err != nil {
		return err
	}

	// Print the logs
	fmt.Printf("%s\n", b)
	return nil
}
