package e2e

import (
	"fmt"
	"strings"
	"testing"

	"github.com/weaveworks/ignite/e2e/util"
	"gotest.tools/assert"
)

// runVMLifecycle is a helper for testing the VM lifecycle.
// vmName should be unique for each test.
func runVMLifecycle(t *testing.T, vmName, runtime, networkPlugin string) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	igniteCmd := util.NewCommand(t, igniteBin)

	// Clean-up the following VM.
	defer igniteCmd.New().
		With("rm", "-f", vmName).
		Run()

	// Run VM.
	igniteCmd.New().
		WithRuntime(runtime).
		WithNetwork(networkPlugin).
		With("run").
		With("--name=" + vmName).
		With("--ssh").
		With(util.DefaultVMImage).
		Run()

	// Check network access.
	igniteCmd.New().
		With("exec", vmName).
		With("curl", "google.com").
		Run()

	// Stop VM.
	igniteCmd.New().
		With("stop", vmName).
		Run()

	// Start VM with explicit runtime and network because these info are
	// removed when VM is stopped. Without it, the VM will start with the
	// default provider configurations.
	igniteCmd.New().
		With("start", vmName).
		WithRuntime(runtime).
		WithNetwork(networkPlugin).
		Run()

	// Check network access after reboot.
	igniteCmd.New().
		With("exec", vmName).
		With("curl", "google.com").
		Run()
}

func TestVMLifecycleWithDockerAndDockerBridge(t *testing.T) {
	runVMLifecycle(
		t,
		"e2e-test-vm-lifecycle-docker-and-docker-bridge",
		"docker",
		"docker-bridge",
	)
}

func TestVMLifecycleWithDockerAndCNI(t *testing.T) {
	runVMLifecycle(
		t,
		"e2e-test-vm-lifecycle-docker-and-cni",
		"docker",
		"cni",
	)
}

func TestVMLifecycleWithContainerdAndCNI(t *testing.T) {
	runVMLifecycle(
		t,
		"e2e-test-vm-lifecycle-containerd-and-cni",
		"containerd",
		"cni",
	)
}

// TestVMProviderSwitch tests that a VM's runtime and network-plugin can be
// changed.
func TestVMProviderSwitch(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	vmName := "e2e-test-vm-providers-switch"

	igniteCmd := util.NewCommand(t, igniteBin)

	// Clean-up the following VM.
	defer igniteCmd.New().
		With("rm", "-f", vmName).
		Run()

	// Run VM.
	igniteCmd.New().
		WithRuntime("containerd").
		WithNetwork("cni").
		With("run").
		With("--name=" + vmName).
		With("--ssh").
		With(util.DefaultVMImage).
		Run()

	// Check network access.
	igniteCmd.New().
		With("exec", vmName).
		With("curl", "google.com").
		Run()

	// Stop VM.
	igniteCmd.New().
		With("stop", vmName).
		Run()

	// Start VM with different providers.
	igniteCmd.New().
		With("start", vmName).
		WithRuntime("docker").
		WithNetwork("docker-bridge").
		Run()

	// Check network access.
	igniteCmd.New().
		With("exec", vmName).
		With("curl", "google.com").
		Run()
}

func TestVMStartNonDefaultProvider(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	vmName := "e2e_test_vm_start_non_default_providers"

	igniteCmd := util.NewCommand(t, igniteBin)

	wantInspect := "docker docker-bridge"

	// Clean-up the following VM.
	defer igniteCmd.New().
		With("rm", "-f", vmName).
		Run()

	// Create VM.
	igniteCmd.New().
		WithRuntime("docker").
		WithNetwork("docker-bridge").
		With("create").
		With(util.DefaultVMImage).
		With("--name=" + vmName).
		Run()

	// Start the VM.
	igniteCmd.New().
		With("start", vmName).
		Run()

	// Inspect the VM runtime and network-plugin.
	inspect := igniteCmd.New().
		With("inspect", "vm", vmName).
		With("-t", "{{.Status.Runtime.Name}} {{.Status.Network.Plugin}}")
	inspectOut, inspectErr := inspect.Cmd.CombinedOutput()
	assert.Check(t, inspectErr, fmt.Sprintf("cmd: \n%q\n%s", inspect.Cmd, inspectOut))
	gotInspect := strings.TrimSpace(string(inspectOut))
	assert.Equal(t, gotInspect, wantInspect, fmt.Sprintf("unexpected VM properties:\n\t(WNT): %q\n\t(GOT): %q", wantInspect, gotInspect))
}
