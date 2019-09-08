//
// This is the e2e package to run tests for Ignite
// Currently, we do local e2e tests only
// we have to wait until the CI setup to allow Ignite to run with sudo and in a KVM environment.
//
// How to run tests:
// sudo IGNITE_E2E_HOME=$PWD $(which go) test ./e2e/.
//

package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"

	"gotest.tools/assert"
)

func TestIgniteRunWithDockerAndDockerBridge(t *testing.T) {
	// vmName should be unique for each test
	const vmName = "e2e_test_ignite_run_docker_docker_bridge"

	dir := os.Getenv("IGNITE_E2E_HOME")
	assert.Assert(t, dir != "", "IGNITE_E2E_HOME should be set")

	binary := path.Join(dir, "bin/ignite")
	cmd := exec.Command(binary,
		"--runtime=docker",
		"--network-plugin=docker-bridge",
		"run", "--name=" + vmName,
		"weaveworks/ignite-ubuntu")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	defer func() {
		cmd := exec.Command(binary,
			"--runtime=docker",
			"--network-plugin=docker-bridge",
			"rm", "-f", vmName)
		assert.Check(t, cmd.Run(), "vm removal should not fail")
	}()

	assert.Check(t, err, fmt.Sprintf("%q should not fail to run", cmd.Args))
}
