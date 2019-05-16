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
	"os"
	"path"
)

type state int

const (
	Stopped state = iota
	Running
)

type vmMetadata struct {
	ID      string `json:"ID"`
	Name    string `json:"Name"`
	ImageID string `json:"ImageID"`
	State   state  `json:"State"`
}

// NewCmdCreate creates a new VM from an image
func NewCmdCreate(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [image] [Name]",
		Short: "Create a new containerized VM without starting it",
		Args:  cobra.MinimumNArgs(2),
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

	// Check if given image exists TODO: Selection by Name
	if dir, err := os.Stat(path.Join(constants.IMAGE_DIR, imageID)); os.IsNotExist(err) || !dir.IsDir() {
		return fmt.Errorf("not an image: %s", imageID)
	}

	// Create a new ID for the VM
	vmID, err := util.NewID(constants.VM_DIR)
	if err != nil {
		return err
	}

	md := &vmMetadata{
		ID:      vmID,
		Name:    args[1],
		ImageID: imageID,
		State:   Stopped,
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

func (md vmMetadata) save() error {
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

func (md vmMetadata) copyImage() error {
	if err := util.CopyFile(path.Join(constants.IMAGE_DIR, md.ImageID, constants.IMAGE_FS),
		path.Join(constants.VM_DIR, md.ID, constants.IMAGE_FS)); err != nil {
		return errors.Wrapf(err, "failed to copy image %s to VM %s", md.ImageID, md.ID)
	}

	return nil
}
