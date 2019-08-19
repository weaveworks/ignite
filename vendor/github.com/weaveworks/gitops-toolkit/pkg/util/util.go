package util

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func ExecuteCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command %q exited with %q: %v", cmd.Args, out, err)
	}

	return string(bytes.TrimSpace(out)), nil
}

func MatchPrefix(prefix string, fields ...string) ([]string, bool) {
	var prefixMatches, exactMatches []string

	for _, str := range fields {
		if str == prefix {
			exactMatches = append(exactMatches, str)
		} else if strings.HasPrefix(str, prefix) {
			prefixMatches = append(prefixMatches, str)
		}
	}

	// If we have exact matches, return them
	// and set the exact match boolean
	if len(exactMatches) > 0 {
		return exactMatches, true
	}

	return prefixMatches, false
}
