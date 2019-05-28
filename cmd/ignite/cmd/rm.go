package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
	"os"
)

type rmOptions struct {
	vm    *vmmd.VMMetadata
	force bool
}

// NewCmdRm removes the given VM
func NewCmdRm(out io.Writer) *cobra.Command {
	ro := &rmOptions{}

	cmd := &cobra.Command{
		Use:   "rm [id]",
		Short: "Remove a VM",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if ro.vm, err = matchSingleVM(args[0]); err != nil {
					return err
				}
				return RunRm(ro)
			}())
		},
	}

	addRmFlags(cmd.Flags(), ro)
	return cmd
}

func addRmFlags(fs *pflag.FlagSet, ro *rmOptions) {
	fs.BoolVarP(&ro.force, "force", "f", false, "Kill VM if running before removal")
}

func RunRm(ro *rmOptions) error {
	// Check if the VM is running
	if ro.vm.Running() {
		// If force is set, kill the VM
		if ro.force {
			if err := RunKill(&killOptions{
				vm: ro.vm,
			}); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("%s is running", ro.vm.ID)
		}
	}

	if err := os.RemoveAll(ro.vm.ObjectPath()); err != nil {
		return fmt.Errorf("unable to remove directory for %s %q: %v", ro.vm.Type, ro.vm.ID, err)
	}

	fmt.Println(ro.vm.ID)
	return nil
}
