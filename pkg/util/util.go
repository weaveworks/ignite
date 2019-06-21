package util

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"

	"github.com/goombaio/namegenerator"
)

func ExecuteCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
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

func TestRoot() (bool, error) {
	u, err := user.Current()
	if err != nil {
		return false, err
	}

	return u.Uid == "0", nil
}

type Prefixer struct {
	prefix    string
	separator string
}

func NewPrefixer() *Prefixer {
	return &Prefixer{
		prefix:    "ignite", // TODO: Remove the dash from IGNITE_PREFIX and use it here
		separator: "-",
	}
}

func (p *Prefixer) Prefix(input ...string) string {
	if len(input) > 0 {
		p.prefix += p.separator + strings.Join(input, p.separator)
	}

	return p.prefix
}
