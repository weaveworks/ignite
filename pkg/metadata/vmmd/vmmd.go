package vmmd

import (
	"encoding/json"
	"net"

	"github.com/luxas/ignite/pkg/metadata"
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
	ImageID   string
	KernelID  string
	State     state
	VCPUs     int64
	Memory    int64
	IPAddrs   []net.IP
	KernelCmd string
}

func NewVMObjectData(imageID, kernelID string, vCPUs, memory int64, kernelCmd string) *VMObjectData {
	return &VMObjectData{
		KernelID:  kernelID,
		ImageID:   imageID,
		State:     Created,
		VCPUs:     vCPUs,
		Memory:    memory,
		KernelCmd: kernelCmd,
	}
}

func NewVMMetadata(id string, name *metadata.Name, od *VMObjectData) *VMMetadata {
	return &VMMetadata{
		Metadata: metadata.NewMetadata(
			id,
			name,
			metadata.VM,
			od),
	}
}

// The md.ObjectData.(*VMObjectData) assert won't panic as this method can only receive *VMMetadata objects
func (md *VMMetadata) VMOD() *VMObjectData {
	return md.ObjectData.(*VMObjectData)
}
