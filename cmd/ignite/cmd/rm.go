package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/metadata"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"github.com/spf13/cobra"
	"io"
	"os"
)

// NewCmdRm removes the given VM
func NewCmdRm(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm [id]",
		Short: "Remove a VM",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := RunRm(out, cmd, args)
			errutils.Check(err)
		},
	}
	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

func RunRm(out io.Writer, cmd *cobra.Command, args []string) error {
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
	if md.Running() {
		return fmt.Errorf("%s is running", md.ID)
	}

	if err := os.RemoveAll(md.ObjectPath()); err != nil {
		return fmt.Errorf("unable to remove directory for %s %q: %v", md.Type, md.ID, err)
	}

	fmt.Println(md.ID)
	return nil
}
