package v1alpha1

import (
	"testing"
)

func TestNewOCIImageRef(t *testing.T) {
	tests := []struct {
		in, out string
		err     bool
	}{
		{
			in:  "weaveworks/ignite-kernel:4.19.47",
			out: "weaveworks/ignite-kernel:4.19.47",
		},
		{
			in:  "centos",
			out: "centos:latest",
		},
		{
			in:  "skjjnfnskj//bs::777",
			err: true,
		},
	}
	for _, rt := range tests {
		actual, err := NewOCIImageRef(rt.in)
		if (err != nil) != rt.err {
			t.Errorf("expected error %t, actual: %v", rt.err, err)
		}
		if actual.String() != rt.out {
			t.Errorf("expected %q, actual: %q", rt.out, actual.String())
		}
	}
}
