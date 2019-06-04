package vmmd

import (
	"encoding/json"
	"fmt"
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/metadata"
	"net"
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
	VCPUs    int64
	Memory   int64
	IPAddrs  []net.IP
}

func NewVMObjectData(imageID, kernelID string, vCPUs, memory int64) *VMObjectData {
	return &VMObjectData{
		KernelID: kernelID,
		ImageID:  imageID,
		State:    Created,
		VCPUs:    vCPUs,
		Memory:   memory,
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

// The md.ObjectData.(*VMObjectData) assert won't panic as this method can only receive *VMMetadata objects
func (md *VMMetadata) VMOD() *VMObjectData {
	return md.ObjectData.(*VMObjectData)
}
