package run

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/weaveworks/libgitops/pkg/runtime"
)

// Update the test output golden files with this flag.
var update = flag.Bool("update", false, "update inspect output golden files")

// Update the golden files with:
//   go test -v github.com/weaveworks/ignite/cmd/ignite/run -run TestInspect -update
func TestInspect(t *testing.T) {
	cases := []struct {
		name         string
		inspectFlags *InspectFlags
		golden       string
		err          bool
	}{
		{
			name:         "json output",
			inspectFlags: &InspectFlags{OutputFormat: "json"},
			golden:       "output/inspect-json.txt",
		},
		{
			name:         "yaml output",
			inspectFlags: &InspectFlags{OutputFormat: "yaml"},
			golden:       "output/inspect-yaml.txt",
		},
		{
			name:         "template formatted output",
			inspectFlags: &InspectFlags{TemplateFormat: "{{.ObjectMeta.Name}} {{.Spec.Image.OCI}}"},
			golden:       "output/inspect-template.txt",
		},
		{
			name:         "unknown output - error",
			inspectFlags: &InspectFlags{OutputFormat: "text"},
			err:          true,
		},
	}

	for _, rt := range cases {
		t.Run(rt.name, func(t *testing.T) {
			vm, err := createTestVM("someVM", "1699b6ba255cde7f")
			if err != nil {
				t.Fatalf("failed to create test vm: %v", err)
			}

			iop := &inspectOptions{InspectFlags: rt.inspectFlags, object: runtime.Object(vm)}

			// Run inspect and capture stdout.
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err = Inspect(iop)
			if (err != nil) != rt.err {
				t.Errorf("expected error %t, actual: %v", rt.err, err)
			}

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			if err != nil {
				t.Fatalf("failed copying to buffer: %v", err)
			}

			// Construct golden file path.
			goldenFilePath := fmt.Sprintf("testdata%c%s", filepath.Separator, rt.golden)

			// Update the golden file if needed.
			if *update {
				t.Log("update inspect golden files")
				if err := ioutil.WriteFile(goldenFilePath, buf.Bytes(), 0644); err != nil {
					t.Fatalf("failed to update inspect golden file: %s: %v", goldenFilePath, err)
				}
			}

			// Check output only when no error is expected.
			if !rt.err {
				// Read golden file.
				wantOutput, err := ioutil.ReadFile(goldenFilePath)
				if err != nil {
					t.Fatalf("failed to read inspect golden file: %s: %v", goldenFilePath, err)
				}

				if !bytes.Equal(buf.Bytes(), wantOutput) {
					t.Errorf("expected output to be:\n%v\ngot output:\n%v", wantOutput, buf.Bytes())
				}
			}
		})
	}
}
