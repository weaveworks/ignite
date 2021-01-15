package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"gotest.tools/assert"

	"github.com/weaveworks/ignite/e2e/util"
)

// runVolume is a helper for testing volume persistence
// vmName should be unique for each test
func runVolume(t *testing.T, vmName, runtime, networkPlugin string) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	// Create a loop device backed by a test-specific file
	volFile := "/tmp/" + vmName + "_vol"

	createDiskCmd := exec.Command(
		"dd",
		"if=/dev/zero",
		"of="+volFile,
		"bs=1M",
		"count=1024",
	)
	createDiskOut, createDiskErr := createDiskCmd.CombinedOutput()
	assert.Check(t, createDiskErr, fmt.Sprintf("create disk: \n%q\n%s", createDiskCmd.Args, createDiskOut))
	if createDiskErr != nil {
		return
	}
	defer func() {
		os.Remove(volFile)
	}()

	mkfsCmd := exec.Command(
		"mkfs.ext4",
		volFile,
	)
	mkfsOut, mkfsErr := mkfsCmd.CombinedOutput()
	assert.Check(t, mkfsErr, fmt.Sprintf("create disk: \n%q\n%s", mkfsCmd.Args, mkfsOut))
	if mkfsErr != nil {
		return
	}

	losetupCmd := exec.Command(
		"losetup",
		"--find",
		"--show",
		volFile,
	)
	losetupOut, losetupErr := losetupCmd.CombinedOutput()
	assert.Check(t, losetupErr, fmt.Sprintf("vm losetup: \n%q\n%s", losetupCmd.Args, losetupOut))
	if losetupErr != nil {
		return
	}

	loopPath := strings.TrimSpace(string(losetupOut))
	defer func() {
		detachLoopCmd := exec.Command(
			"losetup",
			"--detach",
			loopPath,
		)
		detachLoopOut, detachLoopErr := detachLoopCmd.CombinedOutput()
		assert.Check(t, detachLoopErr, fmt.Sprintf("loop detach: \n%q\n%s", detachLoopCmd.Args, detachLoopOut))
	}()

	igniteCmd := util.NewCommand(t, igniteBin)

	// Clean-up the following VM.
	defer igniteCmd.New().
		With("rm", "-f", vmName).
		Run()

	// Run a vm with the loop-device mounted as a volume @ /my-vol
	igniteCmd.New().
		WithRuntime(runtime).
		WithNetwork(networkPlugin).
		With("run").
		With("--debug"). // work around https://github.com/weaveworks/ignite/issues/679 (leaves dead container after)
		With("--name=" + vmName).
		With("--ssh").
		With("--volumes=" + loopPath + ":/my-vol").
		With(util.DefaultVMImage).
		Run()

	// Touch a file in /my-vol
	igniteCmd.New().
		With("exec", vmName).
		With("touch", "/my-vol/hello-world").
		Run()

	// Stop the vm without force.
	igniteCmd.New().
		With("stop", vmName).
		Run()

	secondVMName := vmName + "-2"

	// Clean-up the following VM.
	defer igniteCmd.New().
		With("rm", "-f", secondVMName).
		Run()

	// Start another vm so we can check my-vol
	igniteCmd.New().
		WithRuntime(runtime).
		WithNetwork(networkPlugin).
		With("run").
		With("--debug"). // work around https://github.com/weaveworks/ignite/issues/679 (leaves dead container after)
		With("--name=" + secondVMName).
		With("--ssh").
		With("--volumes=" + loopPath + ":/my-vol").
		With(util.DefaultVMImage).
		Run()

	// Stat the file in /my-vol using the new vm
	igniteCmd.New().
		With("exec", secondVMName).
		With("stat", "/my-vol/hello-world").
		Run()
}

func TestVolumeWithDockerAndDockerBridge(t *testing.T) {
	runVolume(
		t,
		"e2e-test-volume-docker-and-docker-bridge",
		"docker",
		"docker-bridge",
	)
}

func TestVolumeWithDockerAndCNI(t *testing.T) {
	runVolume(
		t,
		"e2e-test-volume-docker-and-cni",
		"docker",
		"cni",
	)
}

func TestVolumeWithContainerdAndCNI(t *testing.T) {
	runVolume(
		t,
		"e2e-test-volume-containerd-and-cni",
		"containerd",
		"cni",
	)
}
