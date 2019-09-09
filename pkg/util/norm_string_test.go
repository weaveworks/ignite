package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToLower(t *testing.T) {
	utests := []struct {
		name   string
		source []string
		dest   []string
	}{
		{
			name:   "Empty",
			source: []string{},
			dest:   []string{},
		},
		{
			name:   "OneElem",
			source: []string{"First"},
			dest:   []string{"first"},
		},
		{
			name:   "TwoElem",
			source: []string{"first", "SeconD"},
			dest:   []string{"first", "second"},
		},
		{
			name:   "TwoElem",
			source: []string{"first", "SeconD"},
			dest:   []string{"first", "second"},
		},
		{
			name:   "ThreeElem",
			source: []string{"fIrst", "SeconD", "third"},
			dest:   []string{"first", "second", "third"},
		},
	}
	for _, utest := range utests {
		t.Run(utest.name, func(t *testing.T) {
			l := ToLower(utest.source)
			assert.Equal(t, utest.dest, l)
		})
	}
}
