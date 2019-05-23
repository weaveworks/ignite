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
	"text/tabwriter"
)

// NewCmdPs lists running Firecracker VMs
func NewCmdPs(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ps",
		Short: "List running Firecracker VMs",
		Run: func(cmd *cobra.Command, args []string) {
			err := RunPs(out, cmd)
			errutils.Check(err)
		},
	}

	return cmd
}

func RunPs(out io.Writer, cmd *cobra.Command) error {
	var mds []*vmmd.VMMetadata

	// Match all VMs using the VMFilter
	// TODO: VMFilter support for running/stopped VMs
	if matches, err := filter.NewFilterer(vmmd.NewVMFilter(""), metadata.VM.Path(), vmmd.LoadVMMetadata); err == nil {
		if all, err := matches.All(); err == nil {
			if mds, err = vmmd.ToVMMetadataAll(all); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}

	// TODO: tabwriter stuff
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprintln(w, "VM ID\tIMAGE\tKERNEL\tSTATE\tNAME")
	for _, md := range mds {
		od := md.ObjectData.(*vmmd.VMObjectData)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", md.ID, od.ImageID, od.KernelID, od.State, md.Name)
	}
	w.Flush()

	return nil
}
