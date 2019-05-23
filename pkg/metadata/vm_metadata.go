package metadata

import (
	"encoding/json"
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/util"
	"io/ioutil"
	"path"
)

type state int

const (
	VMStopped state = iota
	VMRunning
)

var stateLookup = map[state]string{
	VMStopped: "stopped",
	VMRunning: "running",
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
	*Metadata
}

type VMObjectData struct {
	ImageID  string
	KernelID string
	State    state
}

func NewVMMetadata(id, name, imageID, kernelID string) *VMMetadata {
	return &VMMetadata{
		Metadata: &Metadata{
			ID:   id,
			Name: name,
			Type: VM,
			ObjectData: &VMObjectData{
				ImageID:  imageID,
				KernelID: kernelID,
				State:    VMStopped,
			},
		},
	}
}

func LoadVMMetadata(id string) (*VMMetadata, error) {
	md := &VMMetadata{
		Metadata: &Metadata{
			ID:         id,
			Type:       VM,
			ObjectData: &VMObjectData{},
		},
	}

	err := md.Load()
	return md, err
}

func LoadVMMetadataFilterable() ([]filter.Filterable, error) {
	var mds []filter.Filterable

	entries, err := ioutil.ReadDir(VM.Path())
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			md, err := LoadVMMetadata(entry.Name())
			if err != nil {
				return nil, err
			}

			mds = append(mds, md)
		}
	}

	return mds, nil
}

func ToVMMetadata(f filter.Filterable) (*VMMetadata, error) {
	md, ok := f.(*VMMetadata)
	fmt.Printf("%v: %v\n", f, md)
	if !ok {
		return nil, fmt.Errorf("failed to assert Filterable %v to VMMetadata", f)
	}

	return md, nil
}

func (md *VMMetadata) CopyImage() error {
	od := md.ObjectData.(*VMObjectData)

	if err := util.CopyFile(path.Join(constants.IMAGE_DIR, od.ImageID, constants.IMAGE_FS),
		path.Join(md.ObjectPath(), constants.IMAGE_FS)); err != nil {
		return fmt.Errorf("failed to copy image %q to VM %q: %v", od.ImageID, md.ID, err)
	}

	return nil
}

func (md *VMMetadata) SetState(s state) error {
	md.ObjectData.(*VMObjectData).State = s // Won't panic as this can only receive *VMMetadata objects

	if err := md.Save(); err != nil {
		return err
	}

	return nil
}

func (md *VMMetadata) Running() bool {
	return md.ObjectData.(*VMObjectData).State == VMRunning
}

func (md *VMMetadata) GetKernelID() string {
	return md.ObjectData.(*VMObjectData).KernelID
}
