package util

import (
	"strings"
)

func ToLower(a []string) []string {
	b := make([]string, 0, len(a))
	for _, c := range a {
		b = append(b, strings.ToLower(c))
	}
	return b
}
