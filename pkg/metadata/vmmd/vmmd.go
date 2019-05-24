package vmmd

import (
	"encoding/json"
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/metadata"
	"github.com/luxas/ignite/pkg/util"
	"path"
)

type state int

const (
	Created state = iota
	Stopped
	Running
)

var stateLookup = map[state]string{
	Created: "created",
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

func (x state) String() string {
	return stateLookup[x]
}

type VMMetadata struct {
	*metadata.Metadata
}

type VMObjectData struct {
	ImageID  string
	KernelID string
	State    state
}

func NewVMMetadata(id, name, imageID, kernelID string) *VMMetadata {
	return &VMMetadata{
		Metadata: &metadata.Metadata{
			ID:   id,
			Name: name,
			Type: metadata.VM,
			ObjectData: &VMObjectData{
				ImageID:  imageID,
				KernelID: kernelID,
				State:    Created,
			},
		},
	}
}

func ToVMMetadata(f filter.Filterable) (*VMMetadata, error) {
	md, ok := f.(*VMMetadata)
	if !ok {
		return nil, fmt.Errorf("failed to assert Filterable %v to VMMetadata", f)
	}

	return md, nil
}

func ToVMMetadataAll(a []filter.Filterable) ([]*VMMetadata, error) {
	var mds []*VMMetadata

	for _, f := range a {
		if md, err := ToVMMetadata(f); err == nil {
			mds = append(mds, md)
		} else {
			return nil, err
		}
	}

	return mds, nil
}

// The md.ObjectData.(*VMObjectData) assert won't panic as these methods can only receive *VMMetadata objects
func (md *VMMetadata) CopyImage() error {
	od := md.ObjectData.(*VMObjectData)

	if err := util.CopyFile(path.Join(constants.IMAGE_DIR, od.ImageID, constants.IMAGE_FS),
		path.Join(md.ObjectPath(), constants.IMAGE_FS)); err != nil {
		return fmt.Errorf("failed to copy image %q to VM %q: %v", od.ImageID, md.ID, err)
	}

	return nil
}

func (md *VMMetadata) SetState(s state) error {
	md.ObjectData.(*VMObjectData).State = s

	if err := md.Save(); err != nil {
		return err
	}

	return nil
}

func (md *VMMetadata) Running() bool {
	return md.ObjectData.(*VMObjectData).State == Running
}

func (md *VMMetadata) KernelID() string {
	return md.ObjectData.(*VMObjectData).KernelID
}
