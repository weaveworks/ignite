package containerd

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/runtime"

	v2shim "github.com/containerd/containerd/runtime/v2/shim"
	"gotest.tools/assert"
)

var client runtime.Interface

func init() {
	var clienterr error
	client, clienterr = GetContainerdClient()
	if clienterr != nil {
		panic(clienterr)
	}
}

var imageName, _ = meta.NewOCIImageRef("docker.io/library/busybox:latest")

func TestPull(t *testing.T) {
	err := client.PullImage(imageName)
	if err != nil {
		t.Errorf("Error Pulling image: %s", err)
	}
}

func TestInspect(t *testing.T) {
	result, err := client.InspectImage(imageName)
	t.Log(result)
	if err != nil {
		t.Errorf("Error Inspecting image: %s", err)
	}
}

/*func TestExport(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	fmt.Println(tempDir)

	tarCmd := exec.Command("tar", "-x", "-C", tempDir)
	reader, _, err = client.ExportImage(imageName)
	if err != nil {
		t.Fatal("export err:", err)
	}

	tarCmd.Stdin = reader
	if err := tarCmd.Start(); err != nil {
		t.Fatal(err)
	}

	if err := tarCmd.Wait(); err != nil {
		t.Fatal(err)
	}

	if err := reader.Close(); err != nil {
		t.Fatal(err)
	}
	t.Log("done", tempDir)
}*/

func TestRunRemove(t *testing.T) {
	cName := "ignite-test-foo2"
	cID := "test-foo2"
	vmDir := filepath.Join(constants.VM_DIR, cID)

	// TODO: refactor client RunContainer() to take in generic stateDir
	//       remove dependency on VM constants for runtime client
	//       this specific dir is currently required to support resolvconf
	//       ideally, we could pass any tempdir with any permissions here
	assert.NilError(t, os.MkdirAll(vmDir, constants.DATA_DIR_PERM))

	cfg := &runtime.ContainerConfig{
		Cmd: []string{
			"/bin/sh",
			"-c",
			"echo hello",
		},
		Binds: []*runtime.Bind{
			runtime.BindBoth(vmDir),
		},
		Devices: []*runtime.Bind{
			runtime.BindBoth("/dev/kvm"),
		},
		Labels: map[string]string{},
	}

	taskID, err := client.RunContainer(imageName, cfg, cName, cID)
	if err != nil {
		t.Errorf("Error Running Container /w TaskID %q: %s", taskID, err)
	} else {
		t.Logf("TaskID: %q", taskID)
	}

	// TODO: this works around a race where the task is not yet stopped
	//       do this better -- wait on taskID returned above?
	time.Sleep(time.Second / 4)

	err = client.RemoveContainer(cName)
	if err != nil {
		t.Errorf("Error Removing Container: %s", err)
	}

	// just in case the process is hung -- cleanup
	client.KillContainer(cName, "SIGQUIT") //nolint:errcheck // TODO: common constant for SIGQUIT
	client.RemoveContainer(cName)          //nolint:errcheck
}

func TestV2ShimRuntimesHaveBinaryNames(t *testing.T) {
	for _, runtime := range v2ShimRuntimes {
		if v2shim.BinaryName(runtime) == "" {
			t.Errorf("shim binary could not be found -- %q is an invalid runtime/v2/shim", runtime)
		}
	}
}
