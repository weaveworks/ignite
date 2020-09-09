package run

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/pkg/util"
	"github.com/weaveworks/ignite/pkg/version"
	"sigs.k8s.io/yaml"
)

// versionData provides the version information of ignite.
type versionData struct {
	Ignite      version.Info `json:"igniteVersion"`
	Firecracker version.Info `json:"firecrackerVersion"`
}

// NewCmdVersion provides the version information of ignite
func NewCmdVersion(out io.Writer) *cobra.Command {
	var output string
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of ignite",
		Run: func(cmd *cobra.Command, args []string) {
			util.GenericCheckErr(RunVersion(out, output))
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", output, "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

// RunVersion provides the version information of ignite for the specified format
// TODO: Maybe split this out to a pkg/version/flag package that doesn't import cobra
// for use with ignite-spawn
func RunVersion(out io.Writer, output string) error {
	v := versionData{
		Ignite:      version.GetIgnite(),
		Firecracker: version.GetFirecracker(),
	}

	switch output {
	case "":
		fmt.Fprintf(out, "Ignite version: %#v\n", v.Ignite)
		fmt.Fprintf(out, "Firecracker version: %s\n", v.Firecracker.String())
	case "short":
		fmt.Fprintf(out, "%s\n", v.Ignite.GitVersion)
	case "yaml":
		y, err := yaml.Marshal(&v)
		if err != nil {
			return err
		}
		fmt.Fprintln(out, string(y))
	case "json":
		y, err := json.MarshalIndent(&v, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(out, string(y))
	default:
		return fmt.Errorf("invalid output format: %s", output)
	}

	return nil
}
