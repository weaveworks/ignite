package e2e

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/weaveworks/ignite/e2e/util"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	igniteConstants "github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/dmlegacy"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/network"
	"github.com/weaveworks/ignite/pkg/operations"
	"github.com/weaveworks/ignite/pkg/providers"
	igniteDocker "github.com/weaveworks/ignite/pkg/providers/docker"
	"github.com/weaveworks/ignite/pkg/providers/ignite"
	"github.com/weaveworks/ignite/pkg/runtime"
	igniteUtil "github.com/weaveworks/ignite/pkg/util"
	"gotest.tools/assert"
)

var (
	multinetVM  = "e2e-test-vm-multinet"
	sanboxImage = "weaveworks/ignite:dev"
	kernelImage = "weaveworks/ignite-kernel:5.4.108"
	vmImage     = "weaveworks/ignite-ubuntu"
)

func startAsyncVM(t *testing.T, intfs []string) (*operations.VMChannels, string) {

	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")
	igniteUtil.GenericCheckErr(providers.Populate(ignite.Preload))

	igniteDocker.SetDockerRuntime()
	igniteDocker.SetDockerNetwork()

	providers.RuntimeName = runtime.RuntimeDocker
	providers.NetworkPluginName = network.PluginDockerBridge
	providers.Populate(ignite.Providers)

	vm := providers.Client.VMs().New()
	vm.Status.Runtime.Name = runtime.RuntimeDocker
	vm.Status.Network.Plugin = network.PluginDockerBridge

	ociRef, err := meta.NewOCIImageRef(sanboxImage)
	if err != nil {
		t.Fatalf("Failed to parse OCI image ref %s: %s", sanboxImage, err)
	}
	vm.Spec.Sandbox.OCI = ociRef

	ociRef, err = meta.NewOCIImageRef(kernelImage)
	if err != nil {
		t.Fatalf("Failed to parse OCI image ref %s: %s", kernelImage, err)
	}
	vm.Spec.Kernel.OCI = ociRef
	k, _ := operations.FindOrImportKernel(providers.Client, ociRef)
	vm.SetKernel(k)

	ociRef, err = meta.NewOCIImageRef(vmImage)
	if err != nil {
		t.Fatalf("Failed to parse OCI image ref %s: %s", vmImage, err)
	}
	img, err := operations.FindOrImportImage(providers.Client, ociRef)
	if err != nil {
		t.Fatalf("Failed to find OCI image ref %s: %s", ociRef, err)
	}
	vm.SetImage(img)

	vm.Name = multinetVM
	vm.Spec.SSH = &api.SSH{Generate: true}

	metadata.SetNameAndUID(vm, providers.Client)

	for _, intf := range intfs {
		vm.SetAnnotation(igniteConstants.IGNITE_INTERFACE_ANNOTATION+intf, "tc-redirect")
	}

	providers.Client.VMs().Set(vm)

	err = dmlegacy.AllocateAndPopulateOverlay(vm)
	if err != nil {
		t.Fatalf("Error AllocateAndPopulateOverlay: %s", err)
	}

	vmChans, err := operations.StartVMNonBlocking(vm, false)
	if err != nil {
		t.Fatalf("failed to start a VM: \n%q\n", err)
	}

	return vmChans, vm.GetUID().String()
}

// TestMultipleInterface tests that a VM's can be configured with more than 1 interface
func TestOneExtraInterface(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	igniteCmd := util.NewCommand(t, igniteBin)
	dockerCmd := util.NewCommand(t, runtime.RuntimeDocker.String())

	vmChans, vmID := startAsyncVM(t, []string{"foo"})

	// Clean-up the following VM.
	defer igniteCmd.New().
		With("rm", "-f", multinetVM).
		Run()

	fooAddr := "aa:ca:e9:12:34:56"
	dockerCmd.New().
		With("exec", fmt.Sprintf("ignite-%s", vmID)).
		With("ip", "link", "add", "foo", "type", "veth").
		Run()

	dockerCmd.New().
		With("exec", fmt.Sprintf("ignite-%s", vmID)).
		With("ip", "link", "set", "foo", "address", fooAddr).
		Run()

	// check that the VM has started before trying exec
	if err := <-vmChans.SpawnFinished; err != nil {
		t.Fatalf("failed to start a VM: \n%q\n", err)
	}

	eth1Addr := igniteCmd.New().
		With("exec", multinetVM).
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
    "ignite.weave.works/interface/foo": "tc-redirect"
    "ignite.weave.works/interface/bar": "tc-redirect"
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
		With("--template={{.ObjectMeta.UID}}")

	idOut, idErr := idCmd.Cmd.CombinedOutput()
	assert.Check(t, idErr, fmt.Sprintf("vm id not found: \n%q\n%s", idCmd.Cmd, idOut))
	vmID := strings.TrimSuffix(string(idOut), "\n")

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
	assert.Check(t, strings.Contains(gotEth1Addr, barAddr), fmt.Sprintf("unexpected address found:\n\t(WNT): %q\n\t(GOT): %q", barAddr, gotEth1Addr))

	eth2Addr := igniteCmd.New().
		With("exec", vmName).
		With("cat", "/sys/class/net/eth2/address")

	foundEth2Addr, _ := eth2Addr.Cmd.CombinedOutput()
	gotEth2Addr := strings.TrimSuffix(string(foundEth2Addr), "\n")
	assert.Check(t, strings.Contains(gotEth2Addr, fooAddr), fmt.Sprintf("unexpected address found:\n\t(WNT): %q\n\t(GOT): %q", fooAddr, gotEth2Addr))

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
    "ignite.weave.works/interface/foo": "tc-redirect"
    "ignite.weave.works/interface/bar": "tc-redirect"
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
		With("--template={{.ObjectMeta.UID}}")

	idOut, idErr := idCmd.Cmd.CombinedOutput()
	assert.Check(t, idErr, fmt.Sprintf("vm id not found: \n%q\n%s", idCmd.Cmd, idOut))
	vmID := strings.TrimSuffix(string(idOut), "\n")

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
	assert.Check(t, strings.Contains(gotEth1Addr, barAddr), fmt.Sprintf("unexpected address found:\n\t(WNT): %q\n\t(GOT): %q", barAddr, gotEth1Addr))

	eth2Addr := igniteCmd.New().
		With("exec", vmName).
		With("cat", "/sys/class/net/eth2/address")

	foundEth2Addr, _ := eth2Addr.Cmd.CombinedOutput()
	gotEth2Addr := strings.TrimSuffix(string(foundEth2Addr), "\n")
	assert.Check(t, strings.Contains(gotEth2Addr, fooAddr), fmt.Sprintf("unexpected address found:\n\t(WNT): %q\n\t(GOT): %q", fooAddr, gotEth2Addr))

	eth3Addr := igniteCmd.New().
		With("exec", vmName).
		With("cat", "/sys/class/net/eth3/address")

	_, foundEth3Err := eth3Addr.Cmd.CombinedOutput()
	assert.Error(t, foundEth3Err, "exit status 1", fmt.Sprintf("unexpected output when looking for eth3 : \n%s", foundEth3Err))

}
