package e2e

import (
	"fmt"
	"strings"
	"testing"

	"gotest.tools/assert"

	"github.com/weaveworks/ignite/e2e/util"
)

func TestVMExecInteractive(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	vmName := "e2e_test_ignite_exec_interactive"

	igniteCmd := util.NewCommand(t, igniteBin)

	defer igniteCmd.New().
		With("rm", "-f", vmName).
		Run()

	igniteCmd.New().
		With("run", "--name="+vmName).
		With(util.DefaultVMImage).
		With("--ssh").
		Run()

	// Pass input data from host and write to a file inside the VM.
	remoteFileName := "afile.txt"
	inputContent := "foooo..."
	input := strings.NewReader(inputContent)

	igniteExec := igniteCmd.New().
		With("exec", vmName).
		With("tee", remoteFileName)
	igniteExec.Cmd.Stdin = input
	igniteExec.Run()

	// Check the file content inside the VM.
	catExec := igniteCmd.New().
		With("exec", vmName).
		With("cat", remoteFileName)

	catOut, catErr := catExec.Cmd.CombinedOutput()
	assert.Check(t, catErr, fmt.Sprintf("cat: \n%q\n%s", catExec.Cmd, catOut))
	assert.Equal(t, string(catOut), inputContent, fmt.Sprintf("unexpected file content on host:\n\t(WNT): %q\n\t(GOT): %q", inputContent, string(catOut)))
}
