package e2e

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

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
	assert.Equal(t, string(catOut), wantFileContent, fmt.Sprintf("unexpected copied file content:\n\t(WNT): %q\n\t(GOT): %q", wantFileContent, string(catOut)))
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
	assert.Equal(t, got, want, fmt.Sprintf("unexpected copied file content:\n\t(WNT): %q\n\t(GOT): %q", want, got))
}

func TestCopyDirectoryFromHostToVM(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	// Create a temporary directory on host.
	dir, err := ioutil.TempDir("", "ignite-cp-dir-test")
	if err != nil {
		t.Fatalf("failed to create a directory: %v", err)
	}
	defer os.RemoveAll(dir)

	// Create a file in the directory.
	file, err := ioutil.TempFile(dir, "ignite-cp-file")
	if err != nil {
		t.Fatalf("failed to create a file: %v", err)
	}
	content := []byte("some file content")
	if _, err := file.Write(content); err != nil {
		t.Fatalf("failed to write to file: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Errorf("failed to close file: %v", err)
	}

	vmName := "e2e_test_copy_dir_to_vm"
	source := dir
	dest := fmt.Sprintf("%s:%s", vmName, source)

	// Run a VM.
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

	// Copy dir to VM.
	copyCmd := exec.Command(
		igniteBin,
		"cp", source, dest,
	)
	copyOut, copyErr := copyCmd.CombinedOutput()
	assert.Check(t, copyErr, fmt.Sprintf("copy: \n%q\n%s", copyCmd.Args, copyOut))

	// Check if the directory exists in the VM.
	dirFind := fmt.Sprintf("find %s -type d -name %s", filepath.Dir(source), filepath.Base(source))
	dirFindCmd := exec.Command(
		igniteBin,
		"exec", vmName,
		dirFind,
	)
	dirFindOut, dirFindErr := dirFindCmd.CombinedOutput()
	assert.Check(t, dirFindErr, fmt.Sprintf("find: \n%q\n%s", dirFindCmd.Args, dirFindOut))
	gotDir := strings.TrimSpace(string(dirFindOut))
	assert.Equal(t, gotDir, dir, fmt.Sprintf("unexpected find directory result: \n\t(WNT): %q\n\t(GOT): %q", dir, gotDir))

	// Check if the file inside the directory in the VM has the same content as
	// on the host.
	catCmd := exec.Command(
		igniteBin,
		"exec", vmName,
		"cat", file.Name(),
	)
	catOut, catErr := catCmd.CombinedOutput()
	assert.Check(t, catErr, fmt.Sprintf("cat: \n%q\n%s", catCmd.Args, catOut))
	gotContent := strings.TrimSpace(string(catOut))
	assert.Equal(t, gotContent, string(content), fmt.Sprintf("unexpected copied file content:\n\t(WNT): %q\n\t(GOT): %q", content, gotContent))
}

func TestCopyDirectoryFromVMToHost(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	vmName := "e2e_test_copy_dir_from_vm_to_host"

	// Run a VM.
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

	// Create directory inside the VM.
	rand.Seed(time.Now().UnixNano())
	dirPath := fmt.Sprintf("/tmp/ignite-cp-dir-test%d", rand.Intn(10000))
	mkdir := fmt.Sprintf("mkdir -p %s", dirPath)
	mkdirCmd := exec.Command(
		igniteBin,
		"exec", vmName,
		mkdir,
	)
	mkdirOut, mkdirErr := mkdirCmd.CombinedOutput()
	assert.Check(t, mkdirErr, fmt.Sprintf("mkdir: \n%q\n%s", mkdirCmd.Args, mkdirOut))

	// Create file inside the directory.
	content := "some content on VM"
	filePath := filepath.Join(dirPath, "ignite-cp-file")
	writeFile := fmt.Sprintf("echo %s > %s", content, filePath)
	writeFileCmd := exec.Command(
		igniteBin,
		"exec", vmName,
		writeFile,
	)
	writeFileOut, writeFileErr := writeFileCmd.CombinedOutput()
	assert.Check(t, writeFileErr, fmt.Sprintf("file write: \n%q\n%s", writeFileCmd.Args, writeFileOut))

	// Copy the file to host.
	copyCmd := exec.Command(
		igniteBin,
		"cp",
		fmt.Sprintf("%s:%s", vmName, dirPath),
		dirPath,
	)
	copyOut, copyErr := copyCmd.CombinedOutput()
	assert.Check(t, copyErr, fmt.Sprintf("copy: \n%q\n%s", copyCmd.Args, copyOut))
	defer os.RemoveAll(dirPath)

	// Find copied directory on host.
	if _, err := os.Stat(dirPath); err != nil {
		assert.Check(t, err, fmt.Sprintf("error while checking if dir %q exists: %v", dirPath, err))
	}

	// Check the content of the file inside the copied directory.
	hostContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Errorf("failed to read host file %q content: %v", filePath, err)
	}
	gotContent := strings.TrimSpace(string(hostContent))
	assert.Equal(t, gotContent, content, fmt.Sprintf("unexpected copied file content:\n\t(WNT): %q\n\t(GOT): %q", content, gotContent))
}
