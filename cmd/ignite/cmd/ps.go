package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/metadata"
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
	// Load all VM metadata as Filterable objects
	mdf, err := metadata.LoadVMMetadataFilterable()
	if err != nil {
		return err
	}

	// TODO: Temporary, this will be a RunningVMFilter
	d, err := filter.NewFilterer(metadata.NewVMFilter("")).All(mdf)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", d)

	// Convert the result Filterable to a VMMetadata
	//md, err := metadata.ToVMMetadata(d)
	//if err != nil {
	//	return err
	//}

	var mds []*metadata.VMMetadata

	// Type assert all to VM metadata
	for _, a := range d {
		if md, err := metadata.ToVMMetadata(a); err == nil {
			fmt.Printf("%v\n", md)
			mds = append(mds, md)
		} else {
			return err
		}
	}

	fmt.Printf("mds: %v\n", mds)

	// TODO: tabwriter stuff
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprintln(w, "VM ID\tIMAGE\tKERNEL\tSTATE\tNAME")
	for _, md := range mds {
		fmt.Printf("%v\n", mds)
		od := md.ObjectData.(*metadata.VMObjectData)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", md.ID, od.ImageID, od.KernelID, od.State, md.Name)
	}
	w.Flush()

	return nil
}
