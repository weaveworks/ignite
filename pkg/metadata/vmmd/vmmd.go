package vmmd

import (
	"encoding/json"

	"github.com/weaveworks/ignite/pkg/format"

	"github.com/weaveworks/ignite/pkg/metadata"
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
	ImageID      *metadata.ID
	KernelID     *metadata.ID
	Size         format.DataSize
	State        state
	VCPUs        int64
	Memory       format.DataSize
	IPAddrs      IPAddrs
	PortMappings PortMappings
	KernelCmd    string
}

func NewVMObjectData(imageID, kernelID *metadata.ID, size format.DataSize, vCPUs int64, memory format.DataSize, kernelCmd string) *VMObjectData {
	return &VMObjectData{
		KernelID:  kernelID,
		ImageID:   imageID,
		Size:      size,
		State:     Created,
		VCPUs:     vCPUs,
		Memory:    memory,
		KernelCmd: kernelCmd,
	}
}

func NewVMMetadata(id *metadata.ID, name *metadata.Name, od *VMObjectData) (*VMMetadata, error) {
	md, err := metadata.NewMetadata(id, name, metadata.VM, od)
	if err != nil {
		return nil, err
	}

	return &VMMetadata{Metadata: md}, nil
}

// The md.ObjectData.(*VMObjectData) assert won't panic as this method can only receive *VMMetadata objects
func (md *VMMetadata) VMOD() *VMObjectData {
	return md.ObjectData.(*VMObjectData)
}
