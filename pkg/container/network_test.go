package container

import (
	"testing"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/constants"
)

func TestParseExtraIntfs(t *testing.T) {
	cases := []struct {
		name        string
		annotations string
		wantIntfs   map[string]struct{}
	}{
		{
			name:      "empty object",
			wantIntfs: make(map[string]struct{}),
		},
		{
			name:        "wrong annotations",
			annotations: ",",
			wantIntfs:   make(map[string]struct{}),
		},
		{
			name:        "one interface",
			annotations: "eth1",
			wantIntfs: map[string]struct{}{
				"eth1": {},
			},
		},
		{
			name:        "many interfaces",
			annotations: "eth1,eth2,,eth5,",
			wantIntfs: map[string]struct{}{
				"eth1": {},
				"eth2": {},
				"eth5": {},
			},
		},
		{
			name:        "many interfaces with mainInterface (eth0)",
			annotations: "eth1,eth2,,eth0,",
			wantIntfs: map[string]struct{}{
				"eth1": {},
				"eth2": {},
			},
		},
	}

	for _, rt := range cases {
		t.Run(rt.name, func(t *testing.T) {
			vm := &api.VM{}
			vm.SetAnnotation(constants.IGNITE_EXTRA_INTFS, rt.annotations)

			parsedIntfs := parseExtraIntfs(vm)

			// Check if we're not missing an interface
			for k := range rt.wantIntfs {
				if _, ok := parsedIntfs[k]; !ok {
					t.Errorf("expected interface %q not found in parsed Interfaces: %q", k, parsedIntfs)
				}
			}

			// Check if we don't have extra interfaces
			for k := range parsedIntfs {
				if _, ok := rt.wantIntfs[k]; !ok {
					t.Errorf("found interface %q which was not expected", k)
				}
			}

		})
	}
}
