package run

import (
	"fmt"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"github.com/luxas/ignite/pkg/util"
)

type LogsOptions struct {
	VM *vmmd.VMMetadata
}

func Logs(lo *LogsOptions) error {
	// Check if the VM is running
	if !lo.VM.Running() {
		return fmt.Errorf("%s is not running", lo.VM.ID)
	}

	dockerArgs := []string{
		"logs",
		lo.VM.ID,
	}

	// Fetch the VM logs from docker
	output, err := util.ExecuteCommand("docker", dockerArgs...)
	if err != nil {
		return fmt.Errorf("failed to get logs for VM %q: %v", lo.VM.ID, err)
	}

	// Print the ID and the VM logs
	fmt.Println(lo.VM.ID)
	fmt.Println(output)
	return nil
}
