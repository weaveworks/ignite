package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"os"
	"path"
)

type kernelMetadata struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
}

// NewCmdAddKernel adds a new kernel for VM use
func NewCmdAddKernel(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addkernel [path] [name]",
		Short: "Add an uncompressed kernel image for VM use",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			err := RunAddKernel(out, cmd, args)
			errutils.Check(err)
		},
	}

	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

// RunCreate runs when the Create command is invoked
func RunAddKernel(out io.Writer, cmd *cobra.Command, args []string) error {
	p := args[0]

	if !util.FileExists(p) {
		return fmt.Errorf("not a kernel image: %s", p)
	}

	// Create a new ID for the VM
	kernelID, err := util.NewID(constants.KERNEL_DIR)
	if err != nil {
		return err
	}

	md := &kernelMetadata{
		ID:   kernelID,
		Name: args[1],
	}

	// Save the metadata
	if err := md.save(); err != nil {
		return err
	}

	// Perform the image copy
	// TODO: Replace this with overlayfs
	if err := md.importKernel(p); err != nil {
		return err
	}

	fmt.Println(kernelID)

	return nil
}

func (md kernelMetadata) save() error {
	f, err := os.Create(path.Join(constants.KERNEL_DIR, md.ID, constants.METADATA))
	if err != nil {
		return err
	}
	defer f.Close()

	y, err := json.MarshalIndent(&md, "", "    ")
	if err != nil {
		return err
	}

	if _, err := f.Write(append(y, '\n')); err != nil {
		return err
	}

	return nil
}

func (md kernelMetadata) load() error {
	if md.ID == "" {
		return errors.New("cannot load, kernel metadata ID not set")
	}

	kernelDir := path.Join(constants.KERNEL_DIR, md.ID)

	if dir, err := os.Stat(kernelDir); os.IsNotExist(err) || !dir.IsDir() {
		return fmt.Errorf("not a kernel: %s", md.ID)
	}

	mdFile := path.Join(kernelDir, md.ID)

	if !util.FileExists(mdFile) {
		return fmt.Errorf("metadata file missing for kernel: %s", md.ID)
	}

	d, err := ioutil.ReadFile(mdFile)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(d, &md); err != nil {
		return err
	}
	return nil
}

func (md kernelMetadata) importKernel(p string) error {
	if err := util.CopyFile(p, path.Join(constants.KERNEL_DIR, md.ID, constants.KERNEL_FILE)); err != nil {
		return errors.Wrapf(err, "failed to copy kernel file %s to kernel %s", p, md.ID)
	}

	return nil
}
