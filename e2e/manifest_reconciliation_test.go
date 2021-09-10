package e2e

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/weaveworks/ignite/e2e/util"
	"github.com/weaveworks/ignite/pkg/constants"
	"gotest.tools/assert"
)

func TestIgnitedDaemon(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	igniteCmd := util.NewCommand(t, igniteBin)
	ignitedCmd := util.NewCommand(t, ignitedBin)

	vmName := "manifest-vm"

	// Ensure that the manifest directory exists.
	assert.NilError(t, os.MkdirAll(constants.MANIFEST_DIR, 0755))

	// Write the VM manifest in the manifest directory.
	vmConfig := []byte(`---
apiVersion: ignite.weave.works/v1alpha4
kind: VM
metadata:
  name: manifest-vm
  uid: 599615df99804ae9
spec:
  image:
    oci: weaveworks/ignite-ubuntu
  cpus: 1
  diskSize: 3GB
  memory: 800MB
  ssh: true
status:
  running: true
`)
	vmConfigPath := filepath.Join(constants.MANIFEST_DIR, "test-vm.yaml")
	assert.Check(t, ioutil.WriteFile(vmConfigPath, vmConfig, 0644), "failed to write VM config")
	defer func() {
		assert.NilError(t, os.RemoveAll(vmConfigPath))
	}()

	// Start ignited in daemon mode.
	ignitedDaemon := ignitedCmd.With("daemon").Cmd
	assert.Check(t, ignitedDaemon.Start(), fmt.Sprintf("failed to start ignited daemon:\n%q", ignitedDaemon))
	defer func() {
		assert.Check(t, ignitedDaemon.Process.Kill(), "failed to kill ignited daemon")
	}()

	// Wait for ignited to process the manifest.
	time.Sleep(5 * time.Second)

	// Query and check if the VM is created and running.
	wantStatus := "'true'"
	psArgs := []string{
		"--filter={{.ObjectMeta.Name}}=" + vmName,
		"--template='{{.Status.Running}}'",
	}
	psCmd := igniteCmd.New().
		With("ps").
		With(psArgs...)
	psOut, psErr := psCmd.Cmd.CombinedOutput()
	assert.Check(t, psErr, fmt.Sprintf("ps: \n%q\n%s", psCmd.Cmd, psOut))
	got := strings.TrimSpace(string(psOut))
	assert.Equal(t, got, wantStatus)

	// Delete the manifest.
	assert.NilError(t, os.RemoveAll(vmConfigPath))

	// Wait for ignited to process the manifest deletion.
	time.Sleep(5 * time.Second)

	// Query and check if the VM still exists.
	wantStatus = ""
	psCmd = igniteCmd.New().With("ps").With(psArgs...)
	psOut, psErr = psCmd.Cmd.CombinedOutput()
	assert.Check(t, psErr, fmt.Sprintf("ps: \n%q\n%s", psCmd.Cmd, psOut))
	got = strings.TrimSpace(string(psOut))
	assert.Equal(t, got, wantStatus)
}
