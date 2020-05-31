package e2e

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"gotest.tools/assert"

	"github.com/weaveworks/ignite/pkg/constants"
)

func TestConfigFile(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	cases := []struct {
		name             string
		config           []byte
		args             []string
		wantVMProperties string
		err              bool
	}{
		{
			name:   "invalid config",
			config: []byte(``),
			err:    true,
		},
		{
			name: "minimal valid config",
			config: []byte(`---
apiVersion: ignite.weave.works/v1alpha3
kind: Configuration
`),
			wantVMProperties: fmt.Sprintf("'512.0 MB 1 4.0 GB weaveworks/ignite-ubuntu:latest weaveworks/ignite:dev weaveworks/ignite-kernel:%s <nil>'", constants.DEFAULT_KERNEL_IMAGE_TAG),
		},
		{
			name: "custom vm properties",
			config: []byte(`---
apiVersion: ignite.weave.works/v1alpha3
kind: Configuration
metadata:
  name: test-config
spec:
  vm:
    memory: "2GB"
    diskSize: "3GB"
    cpus: 2
    sandbox:
      oci: weaveworks/ignite:dev
    kernel:
      oci: weaveworks/ignite-kernel:4.19.47
    ssh: true
`),
			wantVMProperties: "'2.0 GB 2 3.0 GB weaveworks/ignite-ubuntu:latest weaveworks/ignite:dev weaveworks/ignite-kernel:4.19.47 {true }'",
		},
		{
			name: "runtime and network config",
			config: []byte(`---
apiVersion: ignite.weave.works/v1alpha3
kind: Configuration
metadata:
  name: test-config
spec:
  runtime: docker
  networkPlugin: docker-bridge
`),
			wantVMProperties: fmt.Sprintf("'512.0 MB 1 4.0 GB weaveworks/ignite-ubuntu:latest weaveworks/ignite:dev weaveworks/ignite-kernel:%s <nil>'", constants.DEFAULT_KERNEL_IMAGE_TAG),
		},
		{
			name: "override properties",
			config: []byte(`---
apiVersion: ignite.weave.works/v1alpha3
kind: Configuration
metadata:
  name: test-config
spec:
  vm:
    memory: "2GB"
    diskSize: "3GB"
    cpus: 2
`),
			args:             []string{"--memory=1GB", "--size=1GB", "--cpus=1", "--ssh"},
			wantVMProperties: fmt.Sprintf("'1024.0 MB 1 1024.0 MB weaveworks/ignite-ubuntu:latest weaveworks/ignite:dev weaveworks/ignite-kernel:%s {true }'", constants.DEFAULT_KERNEL_IMAGE_TAG),
		},
	}

	for _, rt := range cases {
		t.Run(rt.name, func(t *testing.T) {
			// Create config file.
			file, err := ioutil.TempFile("", "ignite-config-file-test")
			if err != nil {
				t.Fatalf("failed to create a file: %v", err)
			}
			defer os.Remove(file.Name())

			// Populate the file.
			_, err = file.Write(rt.config)
			assert.NilError(t, err)
			assert.NilError(t, file.Close())

			vmName := "e2e_test_ignite_config_file"

			// Create a VM with the ignite config file.
			// NOTE: Set a sandbox-image to have deterministic results.
			runArgs := []string{
				"run", "--name=" + vmName,
				"weaveworks/ignite-ubuntu",
				"--ignite-config=" + file.Name(),
				"--sandbox-image=weaveworks/ignite:dev",
			}
			// Append the args to the run args for override flags.
			runArgs = append(runArgs, rt.args...)
			runCmd := exec.Command(
				igniteBin,
				runArgs...,
			)
			runOut, _ := runCmd.CombinedOutput()

			// Check if the VM creation failed with a fatal log message.
			// NOTE: CombinedOutput doesn't return error if the process failed
			// with a fatal message.
			fatalRun := strings.Contains(string(runOut), "level=fatal")

			if !fatalRun {
				// Delete the VM only when the creation succeeds, with the
				// config file.
				defer func() {
					rmvCmd := exec.Command(
						igniteBin,
						"rm", "-f", vmName,
						"--ignite-config="+file.Name(),
					)

					rmvOut, rmvErr := rmvCmd.CombinedOutput()
					assert.Check(t, rmvErr, fmt.Sprintf("vm removal: \n%q\n%s", rmvCmd.Args, rmvOut))
				}()

				// Check if run failure was expected.
				if !fatalRun != rt.err {
					assert.Assert(t, !fatalRun != rt.err, "expected VM creation failure")
				}
			}

			if !rt.err {
				// Query VM properties.
				psCmd := exec.Command(
					igniteBin,
					"ps",
					"--filter={{.ObjectMeta.Name}}="+vmName,
					"--template='{{.Spec.Memory}} {{.Spec.CPUs}} {{.Spec.DiskSize}} {{.Spec.Image.OCI}} {{.Spec.Sandbox.OCI}} {{.Spec.Kernel.OCI}} {{.Spec.SSH}}'",
				)
				psOut, psErr := psCmd.CombinedOutput()
				assert.Check(t, psErr, fmt.Sprintf("ps: \n%q\n%s", psCmd.Args, psOut))
				got := strings.TrimSpace(string(psOut))
				assert.Equal(t, got, rt.wantVMProperties, fmt.Sprintf("unexpected VM properties:\n\t(WNT): %q\n\t(GOT): %q", rt.wantVMProperties, got))
			}
		})
	}
}