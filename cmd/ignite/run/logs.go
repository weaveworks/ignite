package run

import (
	"fmt"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/util"
)

type logsOptions struct {
	vm *api.VM
}

func NewLogsOptions(vmMatch string) (lo *logsOptions, err error) {
	lo = &logsOptions{}
	lo.vm, err = getVMForMatch(vmMatch)
	return
}

func Logs(lo *logsOptions) error {
	// Check if the VM is running
	if !lo.vm.Running() {
		return fmt.Errorf("VM %q is not running", lo.vm.GetUID())
	}

	dockerArgs := []string{
		"logs",
		constants.IGNITE_PREFIX + lo.vm.GetUID().String(),
	}

	// Fetch the VM logs from docker
	output, err := util.ExecuteCommand("docker", dockerArgs...)
	if err != nil {
		return fmt.Errorf("failed to get logs for VM %q: %v", lo.vm.GetUID(), err)
	}

	// Print the ID and the VM logs
	fmt.Println(output)
	return nil
}
