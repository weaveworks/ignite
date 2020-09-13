package e2e

import (
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
		"e2e_test_vm_lifecycle_docker_and_docker_bridge",
		"docker",
		"docker-bridge",
	)
}

func TestVMLifecycleWithDockerAndCNI(t *testing.T) {
	runVMLifecycle(
		t,
		"e2e_test_vm_lifecycle_docker_and_cni",
		"docker",
		"cni",
	)
}

func TestVMLifecycleWithContainerdAndCNI(t *testing.T) {
	runVMLifecycle(
		t,
		"e2e_test_vm_lifecycle_containerd_and_cni",
		"containerd",
		"cni",
	)
}

// TestVMProviderSwitch tests that a VM's runtime and network-plugin can be
// changed.
func TestVMProviderSwitch(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	vmName := "e2e_test_vm_providers_switch"

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