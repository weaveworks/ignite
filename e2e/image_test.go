package e2e

import (
	"os/exec"
	"testing"

	"gotest.tools/assert"

	"github.com/weaveworks/ignite/e2e/util"
)

func TestImportTinyImage(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	// NOTE: Along with tiny image, this also tests the image import failure
	// when there's no /etc directory in the image filesystem.

	testImage := "hello-world:latest"

	igniteCmd := util.NewCommand(t, igniteBin)

	// Remove if the image already exists.
	// Ignore any remove error.
	_, _ = igniteCmd.New().
		With("image", "rm", testImage).
		Cmd.CombinedOutput()

	// Import the image.
	igniteCmd.New().
		With("image", "import", testImage).
		Run()
}

func TestDockerImportImage(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	testImage := "hello-world:latest"

	igniteCmd := util.NewCommand(t, igniteBin)

	// Remove if the image already exists.
	_, _ = igniteCmd.New().
		With("image", "rm", testImage).
		Cmd.CombinedOutput()

	// Remove image from docker image store if already exists.
	rmvDockerImgCmd := exec.Command(
		"docker",
		"rmi", testImage,
	)
	// Ignore error if the image doesn't exists.
	_, _ = rmvDockerImgCmd.CombinedOutput()

	// Import the image.
	igniteCmd.New().
		WithRuntime("docker").
		With("image", "import", testImage).
		Run()
}
