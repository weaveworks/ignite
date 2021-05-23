package e2e

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"gotest.tools/assert"

	"github.com/weaveworks/ignite/e2e/util"
	"github.com/weaveworks/ignite/pkg/runtime"
)

const (
	testOSImage     = "localhost:5000/weaveworks/ignite-ubuntu:test"
	testKernelImage = "localhost:5000/weaveworks/ignite-kernel:test"
)

// client config with auth info for the registry setup in
// e2e/util/setup-private-registry.sh.
// NOTE: Update the auth token if the credentials in setup-private-registry.sh
// is updated.
const clientConfigContent = `
{
        "auths": {
                "localhost:5000": {
                        "auth": "dGVzdHVzZXI6dGVzdHBhc3N3b3Jk"
                }
        }
}
`

func TestPullFromAuthRegistry(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	// Create a client config directory to use in test.
	ccDir, err := ioutil.TempDir("", "ignite-test")
	assert.NilError(t, err)
	defer os.RemoveAll(ccDir)

	templateConfig := `---
apiVersion: ignite.weave.works/v1alpha4
kind: Configuration
metadata:
  name: test-config
spec:
  clientConfigDir: %s
`
	igniteConfigContent := fmt.Sprintf(templateConfig, ccDir)

	cases := []struct {
		name               string
		runtime            runtime.Name
		configWithAuthPath string
		clientConfigFlag   string
		igniteConfig       string
		wantErr            bool
	}{
		{
			name:    "no auth info - containerd",
			runtime: runtime.RuntimeContainerd,
			wantErr: true,
		},
		{
			name:    "no auth info - docker",
			runtime: runtime.RuntimeDocker,
			wantErr: true,
		},
		{
			name:               "client config flag - containerd",
			runtime:            runtime.RuntimeContainerd,
			configWithAuthPath: ccDir,
			clientConfigFlag:   ccDir,
		},
		{
			name:               "client config flag - docker",
			runtime:            runtime.RuntimeDocker,
			configWithAuthPath: ccDir,
			clientConfigFlag:   ccDir,
		},
		{
			name:               "client config in ignite config - containerd",
			runtime:            runtime.RuntimeContainerd,
			configWithAuthPath: ccDir,
			igniteConfig:       igniteConfigContent,
		},
		{
			name:               "client config in ignite config - docker",
			runtime:            runtime.RuntimeDocker,
			configWithAuthPath: ccDir,
			igniteConfig:       igniteConfigContent,
		},
		// Following sets the client config dir to a location without a valid
		// client config file, although the client config dir in the ignite
		// config is correct, the import fails due to bad configuration by the
		// flag override.
		{
			name:               "flag override client config - containerd",
			runtime:            runtime.RuntimeContainerd,
			configWithAuthPath: ccDir,
			clientConfigFlag:   "/tmp",
			igniteConfig:       igniteConfigContent,
			wantErr:            true,
		},
		{
			name:               "flag override client config - docker",
			runtime:            runtime.RuntimeDocker,
			configWithAuthPath: ccDir,
			clientConfigFlag:   "/tmp",
			igniteConfig:       igniteConfigContent,
			wantErr:            true,
		},
		// Following set the client config dir via flag without any actual
		// client config. Import fails due to missing auth info in the given
		// client config dir.
		{
			name:               "invalid client config - containerd",
			runtime:            runtime.RuntimeContainerd,
			configWithAuthPath: "",
			clientConfigFlag:   ccDir,
			wantErr:            true,
		},
		{
			name:               "invalid client config - docker",
			runtime:            runtime.RuntimeDocker,
			configWithAuthPath: "",
			clientConfigFlag:   ccDir,
			wantErr:            true,
		},
	}

	for _, rt := range cases {
		rt := rt
		t.Run(rt.name, func(t *testing.T) {
			igniteCmd := util.NewCommand(t, igniteBin)

			// Remove images from ignite store and runtime store. Remove
			// individually because an error in deleting one image cancels the
			// whole command.
			// TODO: Improve image rm to not fail completely when there are
			// multiple images and some are not found.
			util.RmiCompletely(testOSImage, igniteCmd, rt.runtime)
			util.RmiCompletely(testKernelImage, igniteCmd, rt.runtime)

			// Write client config if given.
			if len(rt.configWithAuthPath) > 0 {
				// Ensure the directory exists and create a config file in the
				// directory.
				assert.NilError(t, os.MkdirAll(rt.configWithAuthPath, 0755))
				configPath := filepath.Join(rt.configWithAuthPath, "config.json")
				assert.NilError(t, os.WriteFile(configPath, []byte(clientConfigContent), 0600))
				defer os.Remove(configPath)
			}

			// Write ignite config if provided.
			var igniteConfigPath string
			if len(rt.igniteConfig) > 0 {
				igniteFile, err := ioutil.TempFile("", "ignite-config-file-test")
				if err != nil {
					t.Fatalf("failed to create a file: %v", err)
				}
				igniteConfigPath = igniteFile.Name()

				_, err = igniteFile.WriteString(rt.igniteConfig)
				assert.NilError(t, err)
				assert.NilError(t, igniteFile.Close())

				defer os.Remove(igniteFile.Name())
			}

			// Construct the ignite image import command.
			imageImportCmdArgs := []string{"--runtime", rt.runtime.String()}
			if len(rt.clientConfigFlag) > 0 {
				imageImportCmdArgs = append(imageImportCmdArgs, "--client-config-dir", rt.clientConfigFlag)
			}
			if len(igniteConfigPath) > 0 {
				imageImportCmdArgs = append(imageImportCmdArgs, "--ignite-config", igniteConfigPath)
			}

			// Run image import.
			_, importErr := igniteCmd.New().
				With("image", "import", testOSImage).
				With(imageImportCmdArgs...).
				Cmd.CombinedOutput()
			if (importErr != nil) != rt.wantErr {
				t.Error("expected OS image import to fail")
			}

			// Run kernel import.
			_, importErr = igniteCmd.New().
				With("image", "import", testKernelImage).
				With(imageImportCmdArgs...).
				Cmd.CombinedOutput()
			if (importErr != nil) != rt.wantErr {
				t.Error("expected kernel image import to fail")
			}
		})
	}
}
