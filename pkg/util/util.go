package util

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/goombaio/namegenerator"
)

func ExecuteCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	//log.Debugf("Command %q returned %q\n", cmdArgs, out)
	if err != nil {
		return "", fmt.Errorf("command %q exited with %q: %v", cmd.Args, out, err)
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

type IDHandler struct {
	ID      string
	baseDir string
	success bool
}

// Creates a new 8-byte ID and handles directory creation/deletion
func NewID(baseDir string) (*IDHandler, error) {
	var id string
	var idPath string
	var idBytes []byte

	for {
		idBytes = make([]byte, 8)
		if _, err := rand.Read(idBytes); err != nil {
			return nil, fmt.Errorf("failed to generate ID: %v", err)
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
		return nil, fmt.Errorf("failed to create directory for ID %q: %v", id, err)
	}

	// Return the generated ID
	return &IDHandler{id, baseDir, false}, nil
}

func (i *IDHandler) Remove() error {
	// If success has not been confirmed, remove the generated directory
	if !i.success {
		if err := os.RemoveAll(path.Join(i.baseDir, i.ID)); err != nil {
			return fmt.Errorf("failed to remove directory for ID %q: %v", i.ID, err)
		}
	}

	return nil
}

func (i *IDHandler) Success() (string, error) {
	i.success = true
	return i.ID, nil
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

func RandomName() string {
	return namegenerator.NewNameGenerator(time.Now().UTC().UnixNano()).Generate()
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
