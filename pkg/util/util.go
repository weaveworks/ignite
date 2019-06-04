package util

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/goombaio/namegenerator"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

func ExecuteCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	cmdArgs := strings.Join(cmd.Args, " ")
	//log.Debugf("Command %q returned %q\n", cmdArgs, out)
	if err != nil {
		return "", errors.Wrapf(err, "command %q exited with %q", cmdArgs, out)
	}

	return string(bytes.TrimSpace(out)), nil
}

func ExecForeground(command string, args ...string) (int, error) {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	cmdArgs := strings.Join(cmd.Args, " ")

	var cmdErr error
	var exitCode int

	if err != nil {
		cmdErr = fmt.Errorf("external command %q exited with an error: %v", cmdArgs, err)

		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			cmdErr = fmt.Errorf("failed to get exit code for external command %q", cmdArgs)
		}
	}

	return exitCode, cmdErr
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
			return "", fmt.Errorf("failed to generate ID: %v", err)
		}

		// Convert the byte slice to a string literally
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

// Fills the given string slice with unique MAC addresses
func NewMAC(buffer *[]string) error {
	var mac string
	var macBytes []byte

	for {
		if len(*buffer) == cap(*buffer) {
			break
		}

		macBytes = make([]byte, 6)
		if _, err := rand.Read(macBytes); err != nil {
			return fmt.Errorf("failed to generate MAC: %v", err)
		}

		// Set local bit, ensure unicast address
		macBytes[0] = (macBytes[0] | 2) & 0xfe

		// Convert the byte slice to a string literally
		mac = fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", macBytes[0], macBytes[1], macBytes[2], macBytes[3], macBytes[4], macBytes[5])

		// If the generated MAC is unique break the generator loop
		unique := true
		for _, testMac := range *buffer {
			if mac == testMac {
				unique = false
				break
			}
		}

		// Generate a new MAC if it's not unique
		if !unique {
			continue
		}

		*buffer = append(*buffer, mac)
	}

	return nil
}

func NewName(s *string) {
	if *s == "" {
		*s = namegenerator.NewNameGenerator(time.Now().UTC().UnixNano()).Generate()
	}
}

func MatchPrefix(prefix string, fields ...string) []string {
	var prefixMatches, exactMatches []string

	for _, str := range fields {
		if strings.HasPrefix(str, prefix) {
			prefixMatches = append(prefixMatches, str)
		}
		if str == prefix {
			exactMatches = append(exactMatches, str)
		}
	}

	// If we have exact matches, return them
	if len(exactMatches) > 0 {
		return exactMatches
	}

	return prefixMatches
}
