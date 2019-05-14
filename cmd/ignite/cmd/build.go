package cmd

import (
	"github.com/luxas/ignite/pkg/image"
	"io"

	"github.com/luxas/ignite/pkg/errors"
	"github.com/spf13/cobra"
)

// buildOptions specifies the properties of a new VM
type buildOptions struct {
	tarPath string
	VMID    []byte
}

// NewCmdBuild builds a Firecracker VM.
func NewCmdBuild(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build [tar]",
		Short: "Build a Firecracker VM",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := RunBuild(out, cmd, args)
			errors.Check(err)
		},
	}

	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

// RunBuild runs when the Build command is invoked
func RunBuild(out io.Writer, cmd *cobra.Command, args []string) error {
	buildOptions, err := newBuildOptions(cmd, args)
	if err != nil {
		return err
	}

	return nil
}

// newBuildOptions constructs a set of options for new VMs
func newBuildOptions(cmd *cobra.Command, args []string) (*buildOptions, error) {
	newID, err := build.NewVMID()
	if err != nil {
		return nil, err
	}

	return &buildOptions{
		tarPath: args[0], // The tar path is given as the first argument
		VMID:    newID,
	}, nil
}
