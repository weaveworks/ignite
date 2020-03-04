package e2e

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/weaveworks/ignite/cmd/ignite/run"
	"gotest.tools/assert"
)

func runCopyFilesToVM(t *testing.T, vmName, source, destination, wantFileContent string) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	runCmd := exec.Command(
		igniteBin,
		"run", "--name="+vmName,
		"weaveworks/ignite-ubuntu",
		"--ssh",
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
	if runErr != nil {
		return
	}

	copyCmd := exec.Command(
		igniteBin,
		"cp", source, destination,
	)
	copyOut, copyErr := copyCmd.CombinedOutput()
	assert.Check(t, copyErr, fmt.Sprintf("copy: \n%q\n%s", copyCmd.Args, copyOut))

	// When copying to a VM, the file path succeeds the file path separator.
	// Split the destination to obtain VM destination file path.
	dest := strings.Split(destination, run.VMFilePathSeparator)
	catCmd := exec.Command(
		igniteBin,
		"exec", vmName,
		"cat", dest[1],
	)
	catOut, catErr := catCmd.CombinedOutput()
	assert.Check(t, catErr, fmt.Sprintf("cat: \n%q\n%s", catCmd.Args, catOut))
	assert.Equal(t, string(catOut), wantFileContent, fmt.Sprintf("unexpected copied file content:\n\tWNT: %q\n\tGOT: %q", wantFileContent, string(catOut)))
}

func TestCopyFileFromHostToVM(t *testing.T) {
	cases := []struct {
		name    string
		content []byte
	}{
		{
			name:    "file_with_content",
			content: []byte("some example file content"),
		},
		{
			name:    "empty_file",
			content: []byte(""),
		},
	}

	for _, rt := range cases {
		t.Run(rt.name, func(t *testing.T) {
			// Create a file.
			file, err := ioutil.TempFile("", "ignite-cp-test")
			if err != nil {
				t.Fatalf("failed to create a file: %v", err)
			}
			defer os.Remove(file.Name())

			// Populate the file.
			if _, err := file.Write(rt.content); err != nil {
				t.Fatalf("failed to write to a file: %v", err)
			}
			if err := file.Close(); err != nil {
				t.Errorf("failed to close file: %v", err)
			}

			vmName := "e2e_test_copy_to_vm_" + rt.name
			runCopyFilesToVM(
				t,
				vmName,
				file.Name(),
				fmt.Sprintf("%s:%s", vmName, file.Name()),
				string(rt.content),
			)
		})
	}
}

func TestCopySymlinkedFileFromHostToVM(t *testing.T) {
	// Create a file.
	file, err := ioutil.TempFile("", "ignite-symlink-cp-test")
	if err != nil {
		t.Fatalf("failed to create a file: %v", err)
	}
	defer os.Remove(file.Name())

	fileContent := []byte("Some file content.")

	if _, err := file.Write(fileContent); err != nil {
		t.Fatalf("failed to write to a file: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Errorf("failed to close file: %v", err)
	}

	// Create a new file symlinked to the first file.
	newName := fmt.Sprintf("%s-link", file.Name())
	os.Symlink(file.Name(), newName)
	defer os.Remove(newName)

	vmName := "e2e_test_copy_symlink_to_vm"

	runCopyFilesToVM(
		t,
		vmName,
		newName,
		fmt.Sprintf("%s:%s", vmName, newName),
		string(fileContent),
	)
}

func TestCopyFileFromVMToHost(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	vmName := "e2e_test_copy_file_from_vm_to_host"

	runCmd := exec.Command(
		igniteBin,
		"run", "--name="+vmName,
		"weaveworks/ignite-ubuntu",
		"--ssh",
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
	if runErr != nil {
		return
	}

	// File to be copied from VM.
	vmFilePath := "/proc/version"
	catCmd := exec.Command(
		igniteBin,
		"exec", vmName,
		"cat", vmFilePath,
	)
	catOut, catErr := catCmd.CombinedOutput()
	assert.Check(t, catErr, fmt.Sprintf("cat: \n%q\n%s", catCmd.Args, catOut))

	// Host file path.
	hostFilePath := "/tmp/ignite-os-version"
	copyCmd := exec.Command(
		igniteBin,
		"cp",
		fmt.Sprintf("%s:%s", vmName, vmFilePath),
		hostFilePath,
	)
	copyOut, copyErr := copyCmd.CombinedOutput()
	assert.Check(t, copyErr, fmt.Sprintf("copy: \n%q\n%s", copyCmd.Args, copyOut))
	defer os.Remove(hostFilePath)

	hostContent, err := ioutil.ReadFile(hostFilePath)
	if err != nil {
		t.Errorf("failed to read host file content: %v", err)
	}

	// NOTE: Since the output of cat in the VM includes newline with "\r\n" but
	// the content of file on the host has "\n" for newline when read using go,
	// trim the whitespaces and compare the result.
	got := strings.TrimSpace(string(hostContent))
	want := strings.TrimSpace(string(catOut))
	assert.Equal(t, got, want, fmt.Sprintf("unexpected copied file content:\n\tWNT: %q\n\tGOT: %q", want, got))
}
