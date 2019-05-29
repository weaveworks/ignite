package run

import (
	"fmt"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"github.com/luxas/ignite/pkg/util"
)

type KillOptions struct {
	VM *vmmd.VMMetadata
}

func Kill(ko *KillOptions) error {
	// Check if the VM is running
	if !ko.VM.Running() {
		return fmt.Errorf("%s is not running", ko.VM.ID)
	}

	dockerArgs := []string{
		"kill",
		"-s",
		"SIGQUIT",
		ko.VM.ID,
	}

	// Kill the VM in docker
	if _, err := util.ExecuteCommand("docker", dockerArgs...); err != nil {
		return fmt.Errorf("failed to kill container for VM %q: %v", ko.VM.ID, err)
	}

	fmt.Println(ko.VM.ID)
	return nil
}
