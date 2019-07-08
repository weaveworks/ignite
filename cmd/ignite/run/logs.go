package run

import (
	"fmt"

	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/filter"

	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/util"
)

type logsOptions struct {
	vm *vmmd.VM
}

func NewLogsOptions(vmMatch string) (*logsOptions, error) {
	lo := &logsOptions{}

	if vm, err := client.VMs().Find(filter.NewIDNameFilter(vmMatch)); err == nil {
		lo.vm = &vmmd.VM{vm}
	} else {
		return nil, err
	}

	return lo, nil
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
	fmt.Println(lo.vm.GetUID())
	fmt.Println(output)
	return nil
}
