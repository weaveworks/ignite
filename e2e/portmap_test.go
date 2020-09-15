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

	"github.com/weaveworks/ignite/e2e/util"
	"gotest.tools/assert"
)

func TestPortmapCleanup(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	vmName := "e2e_test_ignite_portmap_cleanup"
	mappedPort := 4242

	igniteCmd := util.NewCommand(t, igniteBin)

	defer igniteCmd.New().
		With("rm", "-f", vmName).
		Run()

	igniteCmd.New().
		WithNetwork("cni").
		With("run").
		With("--name=" + vmName).
		With("--ssh").
		With("--ports=" + fmt.Sprintf("%d:%d", mappedPort, mappedPort)).
		With(util.DefaultVMImage).
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

	// Check that the IPtable rules are installed
	grepOut, grepErr := grepIPTables(vmID)
	assert.NilError(t, grepErr, fmt.Sprintf("unable to match iptable rules:\n %s", grepOut))

	igniteCmd.New().
		With("rm", "-f", vmName).
		Run()

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
