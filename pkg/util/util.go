package util

import (
	"crypto/rand"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

func ExecuteCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	cmdArgs := strings.Join(cmd.Args, " ")
	//log.Debugf("Command %q returned %q\n", cmdArgs, out)
	if err != nil {
		return "", errors.Wrapf(err, "command %q exited with %q", cmdArgs, out)
	}

	// TODO: strings.Builder?
	return strings.TrimSpace(string(out)), nil
}

func PathExists(path string) (bool, os.FileInfo) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, info
}

func FileExists(filename string) bool {
	exists, info := PathExists(filename)
	if !exists {
		return false
	}
	return !info.IsDir()
}

func DirExists(dirname string) bool {
	exists, info := PathExists(dirname)
	if !exists {
		return false
	}
	return info.IsDir()
}

func CopyFile(src string, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func IsEmptyString(input string) bool {
	return len(strings.TrimSpace(input)) == 0
}

// Creates a new 8-byte ID and return it as a string
func NewID(baseDir string) (string, error) {
	var id string
	var idPath string
	var idBytes []byte

	for {
		idBytes = make([]byte, 8)
		if _, err := rand.Read(idBytes); err != nil {
			return "", errors.Wrap(err, "failed to generate ID")
		}

		// Convert the byte array to a string literally
		id = fmt.Sprintf("%x", idBytes)

		// If the generated ID is unique break the generator loop
		idPath = path.Join(baseDir, id)
		if exists, _ := PathExists(idPath); !exists {
			break
		}
	}

	// Create the directory for the ID
	if err := os.MkdirAll(idPath, os.ModePerm); err != nil {
		return "", errors.Wrapf(err, "failed to create directory for ID: %s", id)
	}

	// Return the generated ID
	return id, nil
}
