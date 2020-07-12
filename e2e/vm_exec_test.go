package e2e

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"gotest.tools/assert"
)

func TestVMExecInteractive(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	vmName := "e2e_test_ignite_exec_interactive"

	runCmd := exec.Command(
		igniteBin,
		"run", "--name="+vmName,
		"--ssh",
		"weaveworks/ignite-ubuntu",
	)
	runOut, runErr := runCmd.CombinedOutput()

	defer func() {
		rmvCmd := exec.Command(
			igniteBin,
			"rm", "-f", vmName,
		)
		rmvOut, rmvErr := rmvCmd.CombinedOutput()
		assert.Check(t, rmvErr, fmt.Sprintf("vm removal: \n%q\n%s", rmvCmd.Args, rmvOut))
	}()

	assert.Check(t, runErr, fmt.Sprintf("vm run: \n%q\n%s", runCmd.Args, runOut))

	// Pass input data from host and write to a file inside the VM.
	remoteFileName := "afile.txt"
	inputContent := "foooo..."
	input := strings.NewReader(inputContent)

	execCmd := exec.Command(
		igniteBin,
		"exec", vmName,
		"tee", remoteFileName,
	)
	execCmd.Stdin = input

	execOut, execErr := execCmd.CombinedOutput()
	assert.Check(t, execErr, fmt.Sprintf("exec: \n%q\n%s", execCmd.Args, execOut))

	// Check the file content inside the VM.
	catCmd := exec.Command(
		igniteBin,
		"exec", vmName,
		"cat", remoteFileName,
	)
	catOut, catErr := catCmd.CombinedOutput()
	assert.Check(t, catErr, fmt.Sprintf("cat: \n%q\n%s", catCmd.Args, catOut))
	assert.Equal(t, string(catOut), inputContent, fmt.Sprintf("unexpected file content on host:\n\t(WNT): %q\n\t(GOT): %q", inputContent, string(catOut)))
}
