package util

import (
	"github.com/pkg/errors"
	"os"
	"os/exec"
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

func IsEmptyString(input string) bool {
	return len(strings.TrimSpace(input)) == 0
}
