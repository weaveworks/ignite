package run

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/util"
	"github.com/luxas/ignite/pkg/version"
	"os"
	"os/exec"
	"path/filepath"
)

type StartOptions struct {
	AttachOptions
	Interactive bool
}

func Start(so *StartOptions) error {
	// Check if the given VM is already running
	if so.VM.Running() {
		return fmt.Errorf("%s is already running", so.VM.ID)
	}

	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		return err
	}
	igniteBinary, _ := filepath.Abs(path)

	dockerArgs := []string{
		"run",
		"-itd",
		"--rm",
		"--name",
		so.VM.ID,
		fmt.Sprintf("-v=%s:/ignite/ignite", igniteBinary),
		fmt.Sprintf("-v=%s:%s", constants.DATA_DIR, constants.DATA_DIR),
		fmt.Sprintf("--stop-timeout=%d", constants.STOP_TIMEOUT+constants.IGNITE_TIMEOUT),
		"--privileged",
		"--device=/dev/kvm",
		fmt.Sprintf("weaveworks/ignite:%s", version.GetFirecracker()),
		so.VM.ID,
	}

	// Start the VM in docker
	if _, err := util.ExecuteCommand("docker", dockerArgs...); err != nil {
		return fmt.Errorf("failed to start container for VM %q: %v", so.VM.ID, err)
	}

	// If starting interactively, attach after starting
	if so.Interactive {
		if err := Attach(&so.AttachOptions); err != nil {
			return err
		}
	} else {
		// Print the ID of the started VM
		fmt.Println(so.VM.ID)
	}

	return nil
}
