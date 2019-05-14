package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/version"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

// Version provides the version information of kubeadm.
type Version struct {
	Ignite      version.Info `json:"igniteVersion"`
	Firecracker version.Info `json:"firecrackerVersion"`
}

// NewCmdVersion provides the version information of kubeadm.
func NewCmdVersion(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of ignite",
		Run: func(cmd *cobra.Command, args []string) {
			err := RunVersion(out, cmd)
			errutils.Check(err)
		},
	}
	cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

// RunVersion provides the version information of kubeadm in format depending on arguments
// specified in cobra.Command.
func RunVersion(out io.Writer, cmd *cobra.Command) error {
	v := Version{
		Ignite:      version.GetIgnite(),
		Firecracker: version.GetFirecracker(),
	}

	of, _ := cmd.Flags().GetString("output")
	switch of {
	case "":
		fmt.Fprintf(out, "ignite version: %#v\n", v.Ignite)
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
		return fmt.Errorf("invalid output format: %s", of)
	}
	return nil
}
