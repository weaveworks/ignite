package e2e

import (
	"fmt"
	"os"
	"path"
	"strings"
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

func TestRunWithoutName(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	igniteCmd := util.NewCommand(t, igniteBin)

	labelKey := "testName"
	labelVal := t.Name()

	// Create VM with label to be able to identify the created VM without name.
	igniteCmd.New().
		With("run").
		With(util.DefaultVMImage).
		With("--label", fmt.Sprintf("%s=%s", labelKey, labelVal)).
		Run()

	// List the VM with label and get the VM name.
	psCmd := igniteCmd.New().
		With("ps").
		With("--filter={{.ObjectMeta.Labels}}=~" + labelKey + ":" + labelVal).
		With("--template={{.ObjectMeta.Name}}")
	psOut, psErr := psCmd.Cmd.CombinedOutput()
	assert.Check(t, psErr, fmt.Sprintf("ps: \n%q\n%s", psCmd.Cmd, psOut))
	vmName := strings.TrimSpace(string(psOut))

	// Delete the VM.
	igniteCmd.New().
		With("rm", "-f").
		With(vmName).
		Run()
}
