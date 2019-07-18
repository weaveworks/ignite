package run

import (
	"fmt"
	"github.com/weaveworks/ignite/pkg/runtime/docker"
	"io/ioutil"

	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/util"
)

type logsOptions struct {
	vm *vmmd.VM
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

	// Get the Docker client
	dc, err := docker.GetDockerClient()
	if err != nil {
		return err
	}

	// Fetch the VM logs
	rc, err := dc.ContainerLogs(util.NewPrefixer().Prefix(lo.vm.GetUID()))
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
