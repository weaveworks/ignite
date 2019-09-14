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
	"path"
	"testing"
	"time"

	"gotest.tools/assert"
)

var (
	e2eHome   = os.Getenv("IGNITE_E2E_HOME")
	igniteBin = path.Join(e2eHome, "bin/ignite")
)

// runWithRuntimeAndNetworkPlugin is a helper for running a vm then forcing removal
// vmName should be unique for each test
func runWithRuntimeAndNetworkPlugin(t *testing.T, vmName, runtime, networkPlugin string) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	runCmd := exec.Command(
		igniteBin,
		"--runtime="+runtime,
		"--network-plugin="+networkPlugin,
		"run", "--name="+vmName,
		"weaveworks/ignite-ubuntu",
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
}

func TestIgniteRunWithDockerAndDockerBridge(t *testing.T) {
	runWithRuntimeAndNetworkPlugin(
		t,
		"e2e_test_ignite_run_docker_and_docker_bridge",
		"docker",
		"docker-bridge",
	)
}

func TestIgniteRunWithDockerAndCNI(t *testing.T) {
	runWithRuntimeAndNetworkPlugin(
		t,
		"e2e_test_ignite_run_docker_and_cni",
		"docker",
		"cni",
	)
}

func TestIgniteRunWithContainerdAndCNI(t *testing.T) {
	runWithRuntimeAndNetworkPlugin(
		t,
		"e2e_test_ignite_run_containerd_and_cni",
		"containerd",
		"cni",
	)
}

// runCurl is a helper for testing network connectivity
// vmName should be unique for each test
func runCurl(t *testing.T, vmName, runtime, networkPlugin string) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	runCmd := exec.Command(
		igniteBin,
		"--runtime="+runtime,
		"--network-plugin="+networkPlugin,
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

	time.Sleep(2 * time.Second) // TODO(https://github.com/weaveworks/ignite/issues/423): why is this necessary? Can we work to eliminate this?
	curlCmd := exec.Command(
		igniteBin,
		"--runtime="+runtime,
		"--network-plugin="+networkPlugin,
		"exec", vmName,
		"curl", "google.com",
	)
	curlOut, curlErr := curlCmd.CombinedOutput()
	assert.Check(t, curlErr, fmt.Sprintf("curl: \n%q\n%s", curlCmd.Args, curlOut))
}

func TestCurlWithDockerAndDockerBridge(t *testing.T) {
	runCurl(
		t,
		"e2e_test_curl_docker_and_docker_bridge",
		"docker",
		"docker-bridge",
	)
}

func TestCurlWithDockerAndCNI(t *testing.T) {
	runCurl(
		t,
		"e2e_test_curl_docker_and_cni",
		"docker",
		"cni",
	)
}

func TestCurlWithContainerdAndCNI(t *testing.T) {
	runCurl(
		t,
		"e2e_test_curl_containerd_and_cni",
		"containerd",
		"cni",
	)
}
