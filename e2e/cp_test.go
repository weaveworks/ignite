package e2e

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gotest.tools/assert"

	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/e2e/util"
)

func runCopyFilesToVM(t *testing.T, vmName, source, destination, wantFileContent string) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	igniteCmd := util.NewCommand(t, igniteBin)

	// Clean-up the following VM.
	defer igniteCmd.New().
		With("rm", "-f", vmName).
		Run()

	// Run a VM.
	igniteCmd.New().
		With("run").
		With("--name=" + vmName).
		With("--ssh").
		With(util.DefaultVMImage).
		Run()

	igniteCmd.New().
		With("cp", source, destination).
		Run()

	// When copying to a VM, the file path succeeds the file path separator.
	// Split the destination to obtain VM destination file path.
	dest := strings.Split(destination, run.VMFilePathSeparator)

	catCmd := igniteCmd.New().
		With("exec", vmName).
		With("cat", dest[1])

	catOut, catErr := catCmd.Cmd.CombinedOutput()
	assert.Check(t, catErr, fmt.Sprintf("cat: \n%q\n%s", catCmd.Cmd, catOut))
	assert.Equal(t, string(catOut), wantFileContent, fmt.Sprintf("unexpected copied file content:\n\t(WNT): %q\n\t(GOT): %q", wantFileContent, string(catOut)))
}

func TestCopyFileFromHostToVM(t *testing.T) {
	cases := []struct {
		name    string
		content []byte
	}{
		{
			name:    "file-with-content",
			content: []byte("some example file content"),
		},
		{
			name:    "empty-file",
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

			vmName := "e2e-test-copy-to-vm-" + rt.name
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
	if err := os.Symlink(file.Name(), newName); err != nil {
		t.Errorf("failed to create symlink: %v", err)
	}
	defer os.Remove(newName)

	vmName := "e2e-test-copy-symlink-to-vm"

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

	vmName := "e2e-test-copy-file-from-vm-to-host"

	igniteCmd := util.NewCommand(t, igniteBin)

	// Clean-up the following VM.
	defer igniteCmd.New().
		With("rm", "-f", vmName).
		Run()

	// Run a VM.
	igniteCmd.New().
		With("run").
		With("--name=" + vmName).
		With("--ssh").
		With(util.DefaultVMImage).
		Run()

	// File to be copied from VM.
	vmFilePath := "/proc/version"
	catCmd := igniteCmd.New().
		With("exec", vmName).
		With("cat", vmFilePath)
	catOut, catErr := catCmd.Cmd.CombinedOutput()
	assert.Check(t, catErr, fmt.Sprintf("cat: \n%q\n%s", catCmd.Cmd, catOut))

	// Host file path.
	hostFilePath := "/tmp/ignite-os-version"
	vmSource := fmt.Sprintf("%s:%s", vmName, vmFilePath)
	igniteCmd.New().
		With("cp", vmSource, hostFilePath).
		Run()
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

	vmName := "e2e-test-copy-dir-to-vm"
	source := dir
	dest := fmt.Sprintf("%s:%s", vmName, source)

	igniteCmd := util.NewCommand(t, igniteBin)

	// Clean-up the following VM.
	defer igniteCmd.New().
		With("rm", "-f", vmName).
		Run()

	// Run a VM.
	igniteCmd.New().
		With("run").
		With("--name=" + vmName).
		With("--ssh").
		With(util.DefaultVMImage).
		Run()

	igniteCmd.New().
		With("cp", source, dest).
		Run()

	// Check if the directory exists in the VM.
	dirFind := fmt.Sprintf("find %s -type d -name %s", filepath.Dir(source), filepath.Base(source))
	dirFindCmd := igniteCmd.New().
		With("exec", vmName).
		With(dirFind)
	dirFindOut, dirFindErr := dirFindCmd.Cmd.CombinedOutput()
	assert.Check(t, dirFindErr, fmt.Sprintf("find: \n%q\n%s", dirFindCmd.Cmd, dirFindOut))
	gotDir := strings.TrimSpace(string(dirFindOut))
	assert.Equal(t, gotDir, dir, fmt.Sprintf("unexpected find directory result: \n\t(WNT): %q\n\t(GOT): %q", dir, gotDir))

	// Check if the file inside the directory in the VM has the same content as
	// on the host.
	catCmd := igniteCmd.New().
		With("exec", vmName).
		With("cat", file.Name())
	catOut, catErr := catCmd.Cmd.CombinedOutput()
	assert.Check(t, catErr, fmt.Sprintf("cat: \n%q\n%s", catCmd.Cmd, catOut))
	gotContent := strings.TrimSpace(string(catOut))
	assert.Equal(t, gotContent, string(content), fmt.Sprintf("unexpected copied file content:\n\t(WNT): %q\n\t(GOT): %q", content, gotContent))
}

func TestCopyDirectoryFromVMToHost(t *testing.T) {
	assert.Assert(t, e2eHome != "", "IGNITE_E2E_HOME should be set")

	vmName := "e2e-test-copy-dir-from-vm-to-host"

	igniteCmd := util.NewCommand(t, igniteBin)

	// Clean-up the following VM.
	defer igniteCmd.New().
		With("rm", "-f", vmName).
		Run()

	// Run a VM.
	igniteCmd.New().
		With("run").
		With("--name=" + vmName).
		With("--ssh").
		With(util.DefaultVMImage).
		Run()

	// Create directory inside the VM.
	rand.Seed(time.Now().UnixNano())
	dirPath := fmt.Sprintf("/tmp/ignite-cp-dir-test%d", rand.Intn(10000))
	mkdir := fmt.Sprintf("mkdir -p %s", dirPath)
	igniteCmd.New().
		With("exec", vmName).
		With(mkdir).
		Run()

	// Create file inside the directory.
	content := "some content on VM"
	filePath := filepath.Join(dirPath, "ignite-cp-file")
	writeFile := fmt.Sprintf("echo %s > %s", content, filePath)
	igniteCmd.New().
		With("exec", vmName).
		With(writeFile).
		Run()

	// Copy the file to host.
	src := fmt.Sprintf("%s:%s", vmName, dirPath)
	igniteCmd.New().
		With("cp", src, dirPath).
		Run()
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
