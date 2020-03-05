package e2e

import (
	"fmt"
	"os/exec"
	"testing"

	"gotest.tools/assert"
)

func TestImportTinyImage(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	// NOTE: Along with tiny image, this also tests the image import failure
	// when there's no /etc directory in the image filesystem.

	testImage := "hello-world:latest"
	// Remove if the image already exists.
	rmvImgCmd := exec.Command(
		igniteBin,
		"image", "rm", testImage,
	)
	// Ignore error if the image doesn't exists.
	_, _ = rmvImgCmd.CombinedOutput()

	// Import the image.
	importImgCmd := exec.Command(
		igniteBin,
		"image", "import", testImage,
	)
	importImgOut, importImgErr := importImgCmd.CombinedOutput()
	assert.Check(t, importImgErr, fmt.Sprintf("image import: \n%q\n%s", importImgCmd.Args, importImgOut))
}
