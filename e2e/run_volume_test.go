//
// This is the e2e package to run tests for Ignite
// Currently, we do local e2e tests only
// we have to wait until the CI setup to allow Ignite to run with sudo and in a KVM environment.
//
// How to run tests:
// sudo IGNITE_E2E_HOME=$PWD $(which go) test ./e2e/. -v -count 1 -run Test
//

package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"gotest.tools/assert"
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

	// Run a vm with the loop-device mounted as a volume @ /my-vol
	runCmd := exec.Command(
		igniteBin,
		"--runtime="+runtime,
		"--network-plugin="+networkPlugin,
		"--volumes="+loopPath+":/my-vol",
		"run", "--name="+vmName,
		"weaveworks/ignite-ubuntu",
		"--ssh",
	)
	runOut, runErr := runCmd.CombinedOutput()

	defer func() {
		rmvCmd := exec.Command(
			igniteBin,
			"--runtime="+runtime,
			"--network-plugin="+networkPlugin,
			"rm", "-f", vmName,
		)
		rmvOut, rmvErr := rmvCmd.CombinedOutput()
		assert.Check(t, rmvErr, fmt.Sprintf("vm removal: \n%q\n%s", rmvCmd.Args, rmvOut))
	}()

	assert.Check(t, runErr, fmt.Sprintf("vm run: \n%q\n%s", runCmd.Args, runOut))
	if runErr != nil {
		return
	}

	// Touch a file in /my-vol
	touchCmd := exec.Command(
		igniteBin,
		"--runtime="+runtime,
		"--network-plugin="+networkPlugin,
		"exec", vmName,
		"touch", "/my-vol/hello-world",
	)
	touchOut, touchErr := touchCmd.CombinedOutput()
	assert.Check(t, touchErr, fmt.Sprintf("touch: \n%q\n%s", touchCmd.Args, touchOut))
	if touchErr != nil {
		return
	}

	// Stop the vm
	stopCmd := exec.Command(
		igniteBin,
		"--runtime="+runtime,
		"--network-plugin="+networkPlugin,
		"stop", vmName,
	)
	stopOut, stopErr := stopCmd.CombinedOutput()
	assert.Check(t, stopErr, fmt.Sprintf("vm stop: \n%q\n%s", stopCmd.Args, stopOut))
	if stopErr != nil {
		return
	}

	// Start another vm so we can check my-vol
	run2Cmd := exec.Command(
		igniteBin,
		"--runtime="+runtime,
		"--network-plugin="+networkPlugin,
		"--volumes="+loopPath+":/my-vol",
		"run", "--name="+vmName+"_2",
		"weaveworks/ignite-ubuntu",
		"--ssh",
	)
	run2Out, run2Err := run2Cmd.CombinedOutput()

	defer func() {
		rmv2Cmd := exec.Command(
			igniteBin,
			"--runtime="+runtime,
			"--network-plugin="+networkPlugin,
			"rm", "-f", vmName+"_2",
		)
		rmv2Out, rmv2Err := rmv2Cmd.CombinedOutput()
		assert.Check(t, rmv2Err, fmt.Sprintf("vm removal: \n%q\n%s", rmv2Cmd.Args, rmv2Out))
	}()

	assert.Check(t, run2Err, fmt.Sprintf("vm run: \n%q\n%s", run2Cmd.Args, run2Out))
	if run2Err != nil {
		return
	}

	// Stat the file in /my-vol using the new vm
	stat2Cmd := exec.Command(
		igniteBin,
		"--runtime="+runtime,
		"--network-plugin="+networkPlugin,
		"exec", vmName+"_2",
		"stat", "/my-vol/hello-world",
	)
	stat2Out, stat2Err := stat2Cmd.CombinedOutput()
	assert.Check(t, stat2Err, fmt.Sprintf("stat2: \n%q\n%s", stat2Cmd.Args, stat2Out))
}

func TestVolumeWithDockerAndDockerBridge(t *testing.T) {
	runVolume(
		t,
		"e2e_test_volume_docker_and_docker_bridge",
		"docker",
		"docker-bridge",
	)
}

func TestVolumeWithDockerAndCNI(t *testing.T) {
	runVolume(
		t,
		"e2e_test_volume_docker_and_cni",
		"docker",
		"cni",
	)
}

func TestVolumeWithContainerdAndCNI(t *testing.T) {
	runVolume(
		t,
		"e2e_test_volume_containerd_and_cni",
		"containerd",
		"cni",
	)
}
