package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/metadata"
	"github.com/luxas/ignite/pkg/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
	"path"
)

type state int

const (
	Stopped state = iota
	Running
)

var stateLookup = map[state]string{
	Stopped: "stopped",
	Running: "running",
}

func (x state) MarshalJSON() ([]byte, error) {
	return json.Marshal(stateLookup[x])
}

func (x *state) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	for k, v := range stateLookup {
		if v == s {
			*x = k
			break
		}
	}

	return nil
}

type vmMetadata struct {
	*metadata.Metadata
}

type vmObjectData struct {
	ImageID  string
	KernelID string
	State    state
}

//type vmMetadata struct {
//	ID       string `json:"ID"`
//	Name     string `json:"Name"`
//	ImageID  string `json:"ImageID"`
//	KernelID string `json:"KernelID"`
//	State    state  `json:"State"`
//}

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
	// Resolve the image ID
	imageID, err := metadata.MatchObject(args[0], metadata.Image)
	if err != nil {
		return err
	}

	// Resolve the kernel ID
	kernelID, err := metadata.MatchObject(args[1], metadata.Kernel)
	if err != nil {
		return err
	}

	// Create a new ID for the VM
	vmID, err := util.NewID(constants.VM_DIR)
	if err != nil {
		return err
	}

	md := &vmMetadata{
		Metadata: &metadata.Metadata{
			ID:   vmID,
			Name: args[2],
			Type: metadata.VM,
			ObjectData: &vmObjectData{
				ImageID:  imageID,
				KernelID: kernelID,
				State:    Stopped,
			},
		},
	}

	// Save the metadata
	if err := md.Save(); err != nil {
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

func loadVMMetadata(vmID string) (*vmMetadata, error) {
	md := &vmMetadata{
		Metadata: &metadata.Metadata{
			ID:         vmID,
			Type:       metadata.VM,
			ObjectData: &vmObjectData{},
		},
	}

	if err := md.Load(); err != nil {
		return nil, fmt.Errorf("failed to load VM metadata: %v", err)
	}

	return md, nil
}

func (md *vmMetadata) copyImage() error {
	od := md.ObjectData.(*vmObjectData)

	if err := util.CopyFile(path.Join(constants.IMAGE_DIR, od.ImageID, constants.IMAGE_FS),
		path.Join(md.ObjectPath(), constants.IMAGE_FS)); err != nil {
		return errors.Wrapf(err, "failed to copy image %q to VM %q", od.ImageID, md.ID)
	}

	return nil
}

func (md *vmMetadata) setState(s state) error {
	md.ObjectData.(*vmObjectData).State = s // Won't panic as this can only receive *vmMetadata objects

	if err := md.Save(); err != nil {
		return err
	}

	return nil
}

func (md *vmMetadata) running() bool {
	return md.ObjectData.(*vmObjectData).State == Running
}
