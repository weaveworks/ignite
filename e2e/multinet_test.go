package e2e

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/weaveworks/ignite/e2e/util"
	"gotest.tools/assert"
)

// TestMultipleInterface tests that a VM's can be configured with more than 1 interface
func TestOneExtraInterface(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	vmName := "e2e-test-vm-multinet"

	igniteCmd := util.NewCommand(t, igniteBin)
	dockerCmd := util.NewCommand(t, "docker")

	// Clone this repo in a new dir.
	tempDir, err := ioutil.TempDir("", "ignite-multinet")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Write a VM config with annotations

	vmConfig := []byte(`---
apiVersion: ignite.weave.works/v1alpha4
kind: VM
metadata:
  name: e2e-test-vm-multinet
  annotations:
    "ignite.weave.works/extra-intfs": "foo"
spec:
  image:
    oci: weaveworks/ignite-ubuntu
  cpus: 1
  diskSize: 3GB
  memory: 800MB
  ssh: true
`)

	vmConfigPath := filepath.Join(tempDir, "my-vm.yaml")
	assert.Check(t, ioutil.WriteFile(vmConfigPath, vmConfig, 0644), "failed to write VM config")

	// Clean-up the following VM.
	defer igniteCmd.New().
		With("rm", "-f", vmName).
		Run()

	// Run VM.
	igniteCmd.New().
		WithRuntime("docker").
		WithNetwork("docker-bridge").
		With("run").
		With("--ssh").
		With("--wait=false").
		With("--config=" + vmConfigPath).
		Run()

	// Get the VM ID
	idCmd := igniteCmd.New().
		With("ps").
		With("--filter").
		With(fmt.Sprintf("{{.ObjectMeta.Name}}=%s", vmName)).
		With("--quiet")

	idOut, idErr := idCmd.Cmd.CombinedOutput()
	assert.Check(t, idErr, fmt.Sprintf("vm id not found: \n%q\n%s", idCmd.Cmd, idOut))
	vmID := string(idOut[:len(idOut)-2])

	fooAddr := "aa:ca:e9:12:34:56"
	dockerCmd.New().
		With("exec", fmt.Sprintf("ignite-%s", vmID)).
		With("ip", "link", "add", "foo", "type", "veth").
		Run()

	dockerCmd.New().
		With("exec", fmt.Sprintf("ignite-%s", vmID)).
		With("ip", "link", "set", "foo", "address", fooAddr).
		Run()

	eth1Addr := igniteCmd.New().
		With("exec", vmName).
		With("cat", "/sys/class/net/eth1/address")

	foundEth1Addr, _ := eth1Addr.Cmd.CombinedOutput()
	gotEth1Addr := strings.TrimSuffix(string(foundEth1Addr), "\n")
	assert.Check(t, strings.Contains(gotEth1Addr, fooAddr), fmt.Sprintf("unexpected address found:\n\t(WNT): %q\n\t(GOT): %q", fooAddr, gotEth1Addr))

}

func TestMultipleInterface(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	vmName := "e2e-test-vm-multinet"

	igniteCmd := util.NewCommand(t, igniteBin)
	dockerCmd := util.NewCommand(t, "docker")

	// Clone this repo in a new dir.
	tempDir, err := ioutil.TempDir("", "ignite-multinet")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Write a VM config with annotations

	vmConfig := []byte(`---
apiVersion: ignite.weave.works/v1alpha4
kind: VM
metadata:
  name: e2e-test-vm-multinet
  annotations:
    "ignite.weave.works/extra-intfs": "bar,foo"
spec:
  image:
    oci: weaveworks/ignite-ubuntu
  cpus: 1
  diskSize: 3GB
  memory: 800MB
  ssh: true
`)

	vmConfigPath := filepath.Join(tempDir, "my-vm.yaml")
	assert.Check(t, ioutil.WriteFile(vmConfigPath, vmConfig, 0644), "failed to write VM config")

	// Clean-up the following VM.
	defer igniteCmd.New().
		With("rm", "-f", vmName).
		Run()

	// Run VM.
	igniteCmd.New().
		WithRuntime("docker").
		WithNetwork("docker-bridge").
		With("run").
		With("--ssh").
		With("--wait=false").
		With("--config=" + vmConfigPath).
		Run()

	// Get the VM ID
	idCmd := igniteCmd.New().
		With("ps").
		With("--filter").
		With(fmt.Sprintf("{{.ObjectMeta.Name}}=%s", vmName)).
		With("--quiet")

	idOut, idErr := idCmd.Cmd.CombinedOutput()
	assert.Check(t, idErr, fmt.Sprintf("vm id not found: \n%q\n%s", idCmd.Cmd, idOut))
	vmID := string(idOut[:len(idOut)-2])

	fooAddr := "aa:ca:e9:12:34:56"
	dockerCmd.New().
		With("exec", fmt.Sprintf("ignite-%s", vmID)).
		With("ip", "link", "add", "foo", "type", "veth").
		Run()

	dockerCmd.New().
		With("exec", fmt.Sprintf("ignite-%s", vmID)).
		With("ip", "link", "set", "foo", "address", fooAddr).
		Run()

	barAddr := "aa:ca:e9:12:34:78"
	dockerCmd.New().
		With("exec", fmt.Sprintf("ignite-%s", vmID)).
		With("ip", "link", "add", "bar", "type", "veth").
		Run()

	dockerCmd.New().
		With("exec", fmt.Sprintf("ignite-%s", vmID)).
		With("ip", "link", "set", "bar", "address", barAddr).
		Run()

	eth1Addr := igniteCmd.New().
		With("exec", vmName).
		With("cat", "/sys/class/net/eth1/address")

	foundEth1Addr, _ := eth1Addr.Cmd.CombinedOutput()
	gotEth1Addr := strings.TrimSuffix(string(foundEth1Addr), "\n")
	assert.Check(t, strings.Contains(gotEth1Addr, fooAddr), fmt.Sprintf("unexpected address found:\n\t(WNT): %q\n\t(GOT): %q", fooAddr, gotEth1Addr))

	eth2Addr := igniteCmd.New().
		With("exec", vmName).
		With("cat", "/sys/class/net/eth2/address")

	foundEth2Addr, _ := eth2Addr.Cmd.CombinedOutput()
	gotEth2Addr := strings.TrimSuffix(string(foundEth2Addr), "\n")
	assert.Check(t, strings.Contains(gotEth2Addr, barAddr), fmt.Sprintf("unexpected address found:\n\t(WNT): %q\n\t(GOT): %q", barAddr, gotEth2Addr))

}

func TestMultipleInterfaceImplicit(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	vmName := "e2e-test-vm-multinet"

	igniteCmd := util.NewCommand(t, igniteBin)
	dockerCmd := util.NewCommand(t, "docker")

	// Clone this repo in a new dir.
	tempDir, err := ioutil.TempDir("", "ignite-multinet")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Write a VM config with annotations

	vmConfig := []byte(`---
apiVersion: ignite.weave.works/v1alpha4
kind: VM
metadata:
  name: e2e-test-vm-multinet
  annotations:
    "ignite.weave.works/extra-intfs": "bar,foo"
spec:
  image:
    oci: weaveworks/ignite-ubuntu
  cpus: 1
  diskSize: 3GB
  memory: 800MB
  ssh: true
`)

	vmConfigPath := filepath.Join(tempDir, "my-vm.yaml")
	assert.Check(t, ioutil.WriteFile(vmConfigPath, vmConfig, 0644), "failed to write VM config")

	// Clean-up the following VM.
	defer igniteCmd.New().
		With("rm", "-f", vmName).
		Run()

	// Run VM.
	igniteCmd.New().
		WithRuntime("docker").
		WithNetwork("docker-bridge").
		With("run").
		With("--ssh").
		With("--wait=false").
		With("--config=" + vmConfigPath).
		Run()

	// Get the VM ID
	idCmd := igniteCmd.New().
		With("ps").
		With("--filter").
		With(fmt.Sprintf("{{.ObjectMeta.Name}}=%s", vmName)).
		With("--quiet")

	idOut, idErr := idCmd.Cmd.CombinedOutput()
	assert.Check(t, idErr, fmt.Sprintf("vm id not found: \n%q\n%s", idCmd.Cmd, idOut))
	vmID := string(idOut[:len(idOut)-2])

	fooAddr := "aa:ca:e9:12:34:56"
	dockerCmd.New().
		With("exec", fmt.Sprintf("ignite-%s", vmID)).
		With("ip", "link", "add", "foo", "type", "veth").
		Run()

	dockerCmd.New().
		With("exec", fmt.Sprintf("ignite-%s", vmID)).
		With("ip", "link", "set", "foo", "address", fooAddr).
		Run()

	barAddr := "aa:ca:e9:12:34:78"
	dockerCmd.New().
		With("exec", fmt.Sprintf("ignite-%s", vmID)).
		With("ip", "link", "add", "bar", "type", "veth").
		Run()

	dockerCmd.New().
		With("exec", fmt.Sprintf("ignite-%s", vmID)).
		With("ip", "link", "set", "bar", "address", barAddr).
		Run()

	// this interface should never be found inside a VM
	bazAddr := "aa:ca:e9:12:34:90"
	dockerCmd.New().
		With("exec", fmt.Sprintf("ignite-%s", vmID)).
		With("ip", "link", "add", "baz", "type", "veth").
		Run()

	dockerCmd.New().
		With("exec", fmt.Sprintf("ignite-%s", vmID)).
		With("ip", "link", "set", "baz", "address", bazAddr).
		Run()

	eth1Addr := igniteCmd.New().
		With("exec", vmName).
		With("cat", "/sys/class/net/eth1/address")

	foundEth1Addr, _ := eth1Addr.Cmd.CombinedOutput()
	gotEth1Addr := strings.TrimSuffix(string(foundEth1Addr), "\n")
	assert.Check(t, strings.Contains(gotEth1Addr, fooAddr), fmt.Sprintf("unexpected address found:\n\t(WNT): %q\n\t(GOT): %q", fooAddr, gotEth1Addr))

	eth2Addr := igniteCmd.New().
		With("exec", vmName).
		With("cat", "/sys/class/net/eth2/address")

	foundEth2Addr, _ := eth2Addr.Cmd.CombinedOutput()
	gotEth2Addr := strings.TrimSuffix(string(foundEth2Addr), "\n")
	assert.Check(t, strings.Contains(gotEth2Addr, barAddr), fmt.Sprintf("unexpected address found:\n\t(WNT): %q\n\t(GOT): %q", barAddr, gotEth2Addr))

	eth3Addr := igniteCmd.New().
		With("exec", vmName).
		With("cat", "/sys/class/net/eth3/address")

	_, foundEth3Err := eth3Addr.Cmd.CombinedOutput()
	assert.Error(t, foundEth3Err, "exit status 1", fmt.Sprintf("unexpected output when looking for eth3 : \n%s", foundEth3Err))

}
