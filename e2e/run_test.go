package e2e

import (
	"os"
	"path"
	"testing"

	"gotest.tools/assert"

	"github.com/weaveworks/ignite/e2e/util"
)

var (
	e2eHome   = os.Getenv("IGNITE_E2E_HOME")
	igniteBin = path.Join(e2eHome, "bin/ignite")
)

// runWithRuntimeAndNetworkPlugin is a helper for running a vm then forcing removal
// vmName should be unique for each test
func runWithRuntimeAndNetworkPlugin(t *testing.T, vmName, runtime, networkPlugin string) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	igniteCmd := util.NewCommand(t, igniteBin)

	defer igniteCmd.New().
		With("rm", "-f").
		With(vmName).
		Run()

	igniteCmd.New().
		WithRuntime(runtime).
		WithNetwork(networkPlugin).
		With("run").
		With("--name=" + vmName).
		With(util.DefaultVMImage).
		Run()
}

func TestIgniteRunWithDockerAndDockerBridge(t *testing.T) {
	runWithRuntimeAndNetworkPlugin(
		t,
		"e2e-test-ignite-run-docker-and-docker-bridge",
		"docker",
		"docker-bridge",
	)
}

func TestIgniteRunWithDockerAndCNI(t *testing.T) {
	runWithRuntimeAndNetworkPlugin(
		t,
		"e2e-test-ignite-run-docker-and-cni",
		"docker",
		"cni",
	)
}

func TestIgniteRunWithContainerdAndCNI(t *testing.T) {
	runWithRuntimeAndNetworkPlugin(
		t,
		"e2e-test-ignite-run-containerd-and-cni",
		"containerd",
		"cni",
	)
}

// runCurl is a helper for testing network connectivity
// vmName should be unique for each test
func runCurl(t *testing.T, vmName, runtime, networkPlugin string) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	igniteCmd := util.NewCommand(t, igniteBin)

	defer igniteCmd.New().
		With("rm", "-f", vmName).
		Run()

	igniteCmd.New().
		WithRuntime(runtime).
		WithNetwork(networkPlugin).
		With("run", "--name="+vmName).
		With(util.DefaultVMImage).
		With("--ssh").
		Run()

	igniteCmd.New().
		With("exec", vmName).
		With("curl", "google.com").
		Run()
}

func TestCurlWithDockerAndDockerBridge(t *testing.T) {
	runCurl(
		t,
		"e2e-test-curl-docker-and-docker-bridge",
		"docker",
		"docker-bridge",
	)
}

func TestCurlWithDockerAndCNI(t *testing.T) {
	runCurl(
		t,
		"e2e-test-curl-docker-and-cni",
		"docker",
		"cni",
	)
}

func TestCurlWithContainerdAndCNI(t *testing.T) {
	runCurl(
		t,
		"e2e-test-curl-containerd-and-cni",
		"containerd",
		"cni",
	)
}
