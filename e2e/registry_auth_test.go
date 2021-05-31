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
	"github.com/weaveworks/ignite/pkg/runtime/containerd"
)

const (
	httpTestOSImage      = "127.5.0.1:5080/weaveworks/ignite-ubuntu:test"
	httpTestKernelImage  = "127.5.0.1:5080/weaveworks/ignite-kernel:test"
	httpsTestOSImage     = "127.5.0.1:5443/weaveworks/ignite-ubuntu:test"
	httpsTestKernelImage = "127.5.0.1:5443/weaveworks/ignite-kernel:test"
)

// client config with auth info for the registry setup in
// e2e/util/setup-private-registry.sh.
// NOTE: Update the auth token if the credentials in setup-private-registry.sh
// is updated.
const clientConfigContent = `
{
	"auths": {
		"http://127.5.0.1:5080": {
			"auth": "aHR0cF90ZXN0dXNlcjpodHRwX3Rlc3RwYXNzd29yZA=="
		},
		"https://127.5.0.1:5443": {
			"auth": "aHR0cHNfdGVzdHVzZXI6aHR0cHNfdGVzdHBhc3N3b3Jk"
		}
	}
}
`

func TestPullFromAuthRegistry(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	os.Setenv(containerd.InsecureRegistriesEnvVar, "http://127.5.0.1:5080,https://127.5.0.1:5443")
	defer os.Unsetenv(containerd.InsecureRegistriesEnvVar)

	// Create a client config directory to use in test.
	emptyDir, err := ioutil.TempDir("", "ignite-test")
	assert.NilError(t, err)
	defer os.RemoveAll(emptyDir)

	// Create a client config directory to use in test.
	ccDir, err := ioutil.TempDir("", "ignite-test")
	assert.NilError(t, err)
	defer os.RemoveAll(ccDir)

	// Ensure the directory exists and create a config file in the
	// directory.
	assert.NilError(t, os.MkdirAll(ccDir, 0755))
	configPath := filepath.Join(ccDir, "config.json")
	assert.NilError(t, os.WriteFile(configPath, []byte(clientConfigContent), 0600))
	defer os.Remove(configPath)

	templateConfig := `---
apiVersion: ignite.weave.works/v1alpha4
kind: Configuration
metadata:
  name: test-config
spec:
  clientConfigDir: %s
`
	igniteConfigContent := fmt.Sprintf(templateConfig, ccDir)

	type testCase struct {
		name             string
		runtime          runtime.Name
		clientConfigFlag string
		igniteConfig     string
		wantErr          bool
	}
	cases := []testCase{
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
			name:             "client config flag - containerd",
			runtime:          runtime.RuntimeContainerd,
			clientConfigFlag: ccDir,
		},
		{
			name:             "client config flag - docker",
			runtime:          runtime.RuntimeDocker,
			clientConfigFlag: ccDir,
		},
		{
			name:         "client config in ignite config - containerd",
			runtime:      runtime.RuntimeContainerd,
			igniteConfig: igniteConfigContent,
		},
		{
			name:         "client config in ignite config - docker",
			runtime:      runtime.RuntimeDocker,
			igniteConfig: igniteConfigContent,
		},
		// Following sets the client config dir to a location without a valid
		// client config file, although the client config dir in the ignite
		// config is correct, the import fails due to bad configuration by the
		// flag override.
		{
			name:             "flag override client config - containerd",
			runtime:          runtime.RuntimeContainerd,
			clientConfigFlag: emptyDir,
			igniteConfig:     igniteConfigContent,
			wantErr:          true,
		},
		{
			name:             "flag override client config - docker",
			runtime:          runtime.RuntimeDocker,
			clientConfigFlag: emptyDir,
			igniteConfig:     igniteConfigContent,
			wantErr:          true,
		},
		// Following sets the client config dir via flag without any actual
		// client config. Import fails due to missing auth info in the given
		// client config dir.
		{
			name:             "invalid client config - containerd",
			runtime:          runtime.RuntimeContainerd,
			clientConfigFlag: emptyDir,
			wantErr:          true,
		},
		{
			name:             "invalid client config - docker",
			runtime:          runtime.RuntimeDocker,
			clientConfigFlag: emptyDir,
			wantErr:          true,
		},
	}

	testFunc := func(rt testCase, osImage, kernelImage string) func(t *testing.T) {
		return func(t *testing.T) {
			igniteCmd := util.NewCommand(t, igniteBin)

			// Remove images from ignite store and runtime store. Remove
			// individually because an error in deleting one image cancels the
			// whole command.
			// TODO: Improve image rm to not fail completely when there are
			// multiple images and some are not found.
			util.RmiCompletely(osImage, igniteCmd, rt.runtime)
			util.RmiCompletely(kernelImage, igniteCmd, rt.runtime)

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
				With("image", "import", osImage).
				With(imageImportCmdArgs...).
				Cmd.CombinedOutput()
			if (importErr != nil) != rt.wantErr {
				t.Errorf("expected error %t, actual: %v", rt.wantErr, importErr)
			}

			// Run kernel import.
			_, importErr = igniteCmd.New().
				With("image", "import", kernelImage).
				With(imageImportCmdArgs...).
				Cmd.CombinedOutput()
			if (importErr != nil) != rt.wantErr {
				t.Errorf("expected error %t, actual: %v", rt.wantErr, importErr)
			}
		}
	}

	for _, rt := range cases {
		rt := rt
		t.Run("http_"+rt.name, testFunc(rt, httpTestOSImage, httpTestKernelImage))
	}

	for _, rt := range cases {
		rt := rt
		t.Run("https_"+rt.name, testFunc(rt, httpsTestOSImage, httpsTestKernelImage))
	}
}
