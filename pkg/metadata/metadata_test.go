package metadata

import (
	"testing"

	"github.com/weaveworks/gitops-toolkit/pkg/runtime"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
)

func TestSetLabels(t *testing.T) {
	cases := []struct {
		name       string
		obj        runtime.Object
		labels     []string
		wantLabels map[string]string
		err        bool
	}{
		{
			name: "nil runtime object",
			obj:  nil,
			err:  true,
		},
		{
			name: "valid labels",
			obj:  &api.VM{},
			labels: []string{
				"label1=value1",
				"label2=value2",
				"label3=",
				"label4=value4,label5=value5",
			},
			wantLabels: map[string]string{
				"label1": "value1",
				"label2": "value2",
				"label3": "",
				"label4": "value4,label5=value5",
			},
		},
		{
			name:   "invalid label - key without value",
			obj:    &api.VM{},
			labels: []string{"key1"},
			err:    true,
		},
		{
			name:   "invalid label - empty name",
			obj:    &api.VM{},
			labels: []string{"="},
			err:    true,
		},
	}

	for _, rt := range cases {
		t.Run(rt.name, func(t *testing.T) {
			err := SetLabels(rt.obj, rt.labels)
			if (err != nil) != rt.err {
				t.Errorf("expected error %t, actual: %v", rt.err, err)
			}

			// Check the values of all the labels.
			for k, v := range rt.wantLabels {
				if rt.obj.GetLabel(k) != v {
					t.Errorf("expected label key %q to have value %q, actual: %q", k, v, rt.obj.GetLabel(k))
				}
			}
		})
	}
}
