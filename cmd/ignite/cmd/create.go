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

type state int

const (
	Stopped state = iota
	Running
)

type vmMetadata struct {
	ID       string `json:"ID"`
	Name     string `json:"Name"`
	ImageID  string `json:"ImageID"`
	KernelID string `json:"KernelID"`
	State    state  `json:"State"`
}

// NewCmdCreate creates a new VM from an image
func NewCmdCreate(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [image] [kernel] [name]",
		Short: "Create a new containerized VM without starting it",
		Args:  cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			err := RunCreate(out, cmd, args)
			errutils.Check(err)
		},
	}

	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

// RunCreate runs when the Create command is invoked
func RunCreate(out io.Writer, cmd *cobra.Command, args []string) error {
	imageID := args[0]
	kernelID := args[1]

	// Check if given image exists TODO: Selection by name
	if !util.DirExists(path.Join(constants.IMAGE_DIR, imageID)) {
		return fmt.Errorf("not an image: %s", imageID)
	}

	// Check if given kernel exists TODO: Selection by name
	if !util.DirExists(path.Join(constants.KERNEL_DIR, kernelID)) {
		return fmt.Errorf("not a kernel: %s", kernelID)
	}

	// Create a new ID for the VM
	vmID, err := util.NewID(constants.VM_DIR)
	if err != nil {
		return err
	}

	md := &vmMetadata{
		ID:       vmID,
		Name:     args[2],
		ImageID:  imageID,
		KernelID: kernelID,
		State:    Stopped,
	}

	// Save the metadata
	if err := md.save(); err != nil {
		return err
	}

	// Perform the image copy
	// TODO: Replace this with overlayfs
	if err := md.copyImage(); err != nil {
		return err
	}

	fmt.Println(vmID)

	return nil
}

func (md *vmMetadata) save() error {
	f, err := os.Create(path.Join(constants.VM_DIR, md.ID, constants.METADATA))
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

func (md *vmMetadata) load() error {
	if md.ID == "" {
		return errors.New("cannot load, VM metadata ID not set")
	}

	vmDir := path.Join(constants.VM_DIR, md.ID)

	if dir, err := os.Stat(vmDir); os.IsNotExist(err) || !dir.IsDir() {
		return fmt.Errorf("not a vm: %s", md.ID)
	}

	mdFile := path.Join(vmDir, constants.METADATA)

	if !util.FileExists(mdFile) {
		return fmt.Errorf("metadata file missing for VM: %s", md.ID)
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

func (md *vmMetadata) copyImage() error {
	if err := util.CopyFile(path.Join(constants.IMAGE_DIR, md.ImageID, constants.IMAGE_FS),
		path.Join(constants.VM_DIR, md.ID, constants.IMAGE_FS)); err != nil {
		return errors.Wrapf(err, "failed to copy image %s to VM %s", md.ImageID, md.ID)
	}

	return nil
}

func (md *vmMetadata) setState(s state) error {
	md.State = s

	if err := md.save(); err != nil {
		return err
	}

	return nil
}
