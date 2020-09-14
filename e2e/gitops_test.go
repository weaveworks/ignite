package e2e

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/weaveworks/ignite/e2e/util"
	"gotest.tools/assert"
)

var (
	fileProtocol = "file://"
	ignitedBin   = path.Join(e2eHome, "bin/ignited")
)

func TestRunGitops(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	igniteCmd := util.NewCommand(t, igniteBin)
	ignitedCmd := util.NewCommand(t, ignitedBin)
	gitCmd := util.NewCommand(t, "git")

	vmName := "my-vm"

	// Create dir for git repo.
	gitDir, err := ioutil.TempDir("", "ignite-gitops")
	if err != nil {
		t.Fatalf("failed to create git repo dir: %v", err)
	}
	defer os.RemoveAll(gitDir)

	// Initialize the repo as a bare git repo.
	gitCmd.Dir(gitDir).
		With("init", "--bare").
		Run()

	// Clone this repo in a new dir.
	cloneDir, err := ioutil.TempDir("", "ignite-gitops-clone")
	if err != nil {
		t.Fatalf("failed to create repo clone dir: %v", err)
	}
	defer os.RemoveAll(cloneDir)

	gitRepoURL := fileProtocol + gitDir

	gitCmd.New().
		With("clone", gitRepoURL, cloneDir).
		Run()

	// Set the git account identity in the cloned repo.
	gitCmd.New().
		Dir(cloneDir).
		With("config", "user.name", "test").
		Run()
	gitCmd.New().
		Dir(cloneDir).
		With("config", "user.email", "test@example.com").
		Run()

	// Write a VM config file in the cloned repo, commit and push.
	vmConfig := []byte(`---
apiVersion: ignite.weave.works/v1alpha3
kind: VM
metadata:
  name: my-vm
  uid: 599615df99804ae8
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

	vmConfigPath := filepath.Join(cloneDir, "my-vm.yaml")
	assert.Check(t, ioutil.WriteFile(vmConfigPath, vmConfig, 0644), "failed to write VM config")

	gitCmd.New().
		Dir(cloneDir).
		With("add", vmConfigPath).
		Run()

	gitCmd.New().
		Dir(cloneDir).
		With("commit", "-m", "add vm config").
		Run()

	gitCmd.New().
		Dir(cloneDir).
		With("push").
		Run()

	// Delete the VM at the end.
	defer igniteCmd.New().
		With("rm", "-f").
		With(vmName).
		Run()

	// Start ignited in gitops mode.
	// NOTE: Running ignited in an empty git repo results in a fatal error. Run
	// ignited in repos with at least one commit.
	ignitedGitops := ignitedCmd.With("gitops", gitRepoURL).Cmd
	assert.Check(t, ignitedGitops.Start(), fmt.Sprintf("failed to start ignited gitops:\n%q", ignitedGitops))
	defer func() {
		assert.Check(t, ignitedGitops.Process.Kill(), "failed to kill ignited gitops")
	}()

	// Wait for ignited to detect the changes and act on it.
	time.Sleep(10 * time.Second)

	wantVMProperties := "'800.0 MB 1 3.0 GB weaveworks/ignite-ubuntu:latest {true } true'"

	psArgs := []string{
		"--filter={{.ObjectMeta.Name}}=" + vmName,
		"--template='{{.Spec.Memory}} {{.Spec.CPUs}} {{.Spec.DiskSize}} {{.Spec.Image.OCI}} {{.Spec.SSH}} {{.Status.Running}}'",
	}
	psCmd := igniteCmd.New().
		With("ps").
		With(psArgs...)
	psOut, psErr := psCmd.Cmd.CombinedOutput()
	assert.Check(t, psErr, fmt.Sprintf("ps: \n%q\n%s", psCmd.Cmd, psOut))

	got := strings.TrimSpace(string(psOut))
	assert.Equal(t, got, wantVMProperties, fmt.Sprintf("unexpected VM properties:\n\t(WNT): %q\n\t(GOT): %q", wantVMProperties, got))
}
