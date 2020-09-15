//
// This is the e2e package to run tests for Ignite
// Currently, we do local e2e tests only
// we have to wait until the CI setup to allow Ignite to run with sudo and in a KVM environment.
//
// How to run tests:
// sudo IGNITE_E2E_HOME=$PWD $(which go) test ./e2e/. -v -count 1 -run Test
//

package e2e

import (
	"fmt"
	"os/exec"
	"testing"

	"gotest.tools/assert"
)

func TestPortmapCleanup(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	vmName := "e2e_test_ignite_portmap"
	mappedPort := 4242

	runCmd := exec.Command(
		igniteBin,
		"run", "--name="+vmName,
		"--ssh",
		"--ports",
		fmt.Sprintf("%d:%d", mappedPort, mappedPort),
		"weaveworks/ignite-ubuntu",
	)
	runOut, runErr := runCmd.CombinedOutput()
	assert.Check(t, runErr, fmt.Sprintf("vm run: \n%q\n%s", runCmd.Args, runOut))

	defer func() {
		rmvCmd := exec.Command(
			"sudo",
			igniteBin,
			"rm", "-f", vmName,
		)
		_ = rmvCmd.Run()
	}()

	// Get the VM ID
	idCmd := exec.Command(
		"ignite",
		"ps",
		"--filter",
		fmt.Sprintf("{{.ObjectMeta.Name}}=%s", vmName),
		"--quiet",
	)
	idOut, idErr := idCmd.CombinedOutput()
	assert.Check(t, idErr, fmt.Sprintf("vm id: \n%q\n%s", idCmd.Args, idOut))
	vmID := string(idOut[:len(idOut)-2])

	// Check that the IPtable rules are installed
	grepOut, grepErr := grepIPTables(vmID)
	assert.NilError(t, grepErr, fmt.Sprintf("unable to match iptable rules:\n %s", grepOut))

	rmvCmd := exec.Command(
		igniteBin,
		"rm", "-f", vmName,
	)
	rmvOut, rmvErr := rmvCmd.CombinedOutput()
	assert.Check(t, rmvErr, fmt.Sprintf("vm removal: \n%q\n%s", rmvCmd.Args, rmvOut))

	// Check that the IPtable rules are removed
	grepOut, grepErr = grepIPTables(vmID)
	assert.Error(t, grepErr, "exit status 1", fmt.Sprintf("unexpected output in grep : \n%s", grepOut))
	assert.Equal(t, len(grepOut), 0, fmt.Sprintf("leftover iptable rules detected:\n%s", string(grepOut)))

}

func grepIPTables(search string) ([]byte, error) {
	grepCmd := exec.Command(
		"grep",
		search,
	)
	natCmd := exec.Command(
		"iptables",
		"-t",
		"nat",
		"-nL",
	)
	pipe, _ := natCmd.StdoutPipe()
	defer pipe.Close()

	grepCmd.Stdin = pipe
	_ = natCmd.Start()
	return grepCmd.Output()
}
