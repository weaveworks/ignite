package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/metadata"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"github.com/luxas/ignite/pkg/util"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdStop kills a Firecracker VM
func NewCmdKill(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kill [id]",
		Short: "Kill a running Firecracker VM",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := RunKill(out, cmd, args)
			errutils.Check(err)
		},
	}

	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

func RunKill(out io.Writer, cmd *cobra.Command, args []string) error {
	var md *vmmd.VMMetadata

	// Match a single VM using the VMFilter
	if matches, err := filter.NewFilterer(vmmd.NewVMFilter(args[0]), metadata.VM.Path(), vmmd.LoadVMMetadata); err == nil {
		if filterable, err := matches.Single(); err == nil {
			if md, err = vmmd.ToVMMetadata(filterable); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}

	// Check if the VM is running
	if !md.Running() {
		return fmt.Errorf("%s is not running", md.ID)
	}

	dockerArgs := []string{
		"kill",
		"-s",
		"SIGQUIT",
		md.ID,
	}

	// Kill the VM in docker
	if _, err := util.ExecuteCommand("docker", dockerArgs...); err != nil {
		return fmt.Errorf("failed to kill container for VM %q: %v", md.ID, err)
	}

	fmt.Println(md.ID)
	return nil
}
