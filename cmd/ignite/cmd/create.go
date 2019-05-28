package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/metadata/imgmd"
	"github.com/luxas/ignite/pkg/metadata/kernmd"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"github.com/luxas/ignite/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
)

type createOptions struct {
	image  *imgmd.ImageMetadata
	kernel *kernmd.KernelMetadata
	vm     *vmmd.VMMetadata
	name   string
	cpus   int64
	memory int64
}

// NewCmdCreate creates a new VM from an image
func NewCmdCreate(out io.Writer) *cobra.Command {
	co := &createOptions{}

	cmd := &cobra.Command{
		// TODO: ValidArgs and different metadata loading setup
		Use:   "create [image] [kernel]",
		Short: "Create a new containerized VM without starting it",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if co.image, err = matchSingleImage(args[0]); err != nil {
					return err
				}
				if co.kernel, err = matchSingleKernel(args[1]); err != nil {
					return err
				}
				return RunCreate(co)
			}())
		},
	}

	addCreateFlags(cmd.Flags(), co)
	return cmd
}

func addCreateFlags(fs *pflag.FlagSet, co *createOptions) {
	addNameFlag(fs, &co.name)
	fs.Int64Var(&co.cpus, "cpus", constants.VM_DEFAULT_CPUS, "VM vCPU count, 1 or even numbers between 1 and 32")
	fs.Int64Var(&co.memory, "memory", constants.VM_DEFAULT_MEMORY, "VM RAM in MiB")
}

func RunCreate(co *createOptions) error {
	// Create a new ID for the VM
	vmID, err := util.NewID(constants.VM_DIR)
	if err != nil {
		return err
	}

	// Create a new name for the VM if none is given
	util.NewName(&co.name)

	// Create new metadata for the VM and add to createOptions for further processing
	// This enables the generated VM metadata to pass straight to start and attach via run
	co.vm = vmmd.NewVMMetadata(vmID, co.name, vmmd.NewVMObjectData(co.image.ID, co.kernel.ID, co.cpus, co.memory))

	// Save the metadata
	if err := co.vm.Save(); err != nil {
		return err
	}

	// Perform the image copy
	// TODO: Replace this with overlayfs
	if err := co.vm.CopyImage(); err != nil {
		return err
	}

	// Print the ID of the created VM
	fmt.Println(co.vm.ID)

	return nil
}
