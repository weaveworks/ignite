package e2e

import (
	"strconv"
	"testing"

	"github.com/weaveworks/ignite/e2e/util"
	"gotest.tools/assert"
)

func TestConcurrentVMCreation(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	igniteCmd := util.NewCommand(t, igniteBin)

	numberOfVMs := 4
	vmNames := []string{}
	cmds := []*util.Command{}

	// Create VM names and VM run commands to execute.
	for i := 1; i <= numberOfVMs; i++ {
		name := "e2e-test-concurrent-vm-create-" + strconv.Itoa(i)
		vmNames = append(vmNames, name)
		cmds = append(
			cmds,
			util.NewCommand(t, igniteBin).
				With("run").
				With("--name="+name).
				With("--ssh").
				With(util.DefaultVMImage),
		)
	}

	// Clean-up the VMs.
	defer igniteCmd.New().
		With("rm", "-f").
		With(vmNames...).
		Run()

	// Run VMs.
	for _, cmd := range cmds {
		assert.Check(t, cmd.Cmd.Start(), "failed to run VM")
	}

	// Wait for all the commands to finish.
	for _, cmd := range cmds {
		assert.Check(t, cmd.Cmd.Wait(), "error waiting for the command to finish")
	}
}
