package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/errutils"
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
	// match all VMs by specifying an empty filter
	ds, err := metadata.NewObjectMatcher("").All(metadata.VM, loadVMMetadata)
	if err != nil {
		return err
	}

	var mds []*vmMetadata

	// Type assert all to VM metadata
	for _, d := range ds {
		mds = append(mds, (*d).(*vmMetadata))

		//if md, err := toVMMetadata(d); err == nil {
		//	mds = append(mds, md)
		//} else {
		//	return err
		//}
	}

	// TODO: tabwriter stuff
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprintln(w, "VM ID\tIMAGE\tKERNEL\tSTATE\tNAME")
	for _, md := range mds {
		od := md.ObjectData.(*vmObjectData)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", md.ID, od.ImageID, od.KernelID, od.State, md.Name)
	}
	w.Flush()

	return nil
}
