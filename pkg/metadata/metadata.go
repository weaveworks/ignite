package metadata

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/util"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type ObjectType int

const (
	Image ObjectType = iota + 1 // Reserve 0 for unset
	Kernel
	VM
)

var ObjectTypeLookup = map[ObjectType]string{
	Image:  "image",
	Kernel: "kernel",
	VM:     "VM",
}

func (x ObjectType) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.String())
}

func (x *ObjectType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	for k, v := range ObjectTypeLookup {
		if v == s {
			*x = k
			break
		}
	}

	return nil
}

func (x ObjectType) String() string {
	return ObjectTypeLookup[x]
}

func (x ObjectType) Path() string {
	switch x {
	case Image:
		return constants.IMAGE_DIR
	case Kernel:
		return constants.KERNEL_DIR
	case VM:
		return constants.VM_DIR
	}

	return ""
}

type ObjectData interface{}

type Metadata struct {
	ID         string     `json:"ID"`
	Name       string     `json:"Name"`
	Type       ObjectType `json:"Type"`
	ObjectData `json:"ObjectData"`
}

func LoadMetadata(i string, t ObjectType) (*Metadata, error) {
	md := &Metadata{
		ID:   i,
		Type: t,
	}

	err := md.Load()
	return md, err
}

func LoadMetadataFilterable(t ObjectType) ([]filter.Filterable, error) {
	var mds []filter.Filterable

	entries, err := ioutil.ReadDir(t.Path())
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			md, err := LoadMetadata(entry.Name(), t)
			if err != nil {
				return nil, err
			}

			mds = append(mds, *md)
		}
	}

	return mds, nil
}

func (md *Metadata) ObjectPath() string {
	return path.Join(md.Type.Path(), md.ID)
}

func (md *Metadata) Save() error {
	f, err := os.Create(path.Join(md.ObjectPath(), constants.METADATA))
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

func (md *Metadata) Load() error {
	if md.ID == "" {
		return errors.New("cannot load metadata, ID not set")
	}

	if md.Type == 0 { // Type is unset
		return errors.New("cannot load metadata, Type not set")
	}

	p := md.ObjectPath()

	if !util.DirExists(p) {
		return fmt.Errorf("nonexistent %s: %s", md.Type, md.ID)
	}

	f := path.Join(p, constants.METADATA)

	if !util.FileExists(f) {
		return fmt.Errorf("metadata file missing for %s: %s", md.Type, md.ID)
	}

	d, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(d, &md); err != nil {
		return err
	}

	return nil
}

// TODO: Move to filter
// Compile-time assert to verify interface compatibility
var _ filter.Filter = &IDNameFilter{}

type IDNameFilter struct {
	prefix string
}

func NewIDNameFilter(p string) *IDNameFilter {
	return &IDNameFilter{
		prefix: p,
	}
}

func (n *IDNameFilter) Filter(f filter.Filterable) (bool, error) {
	md := f.(*Metadata)
	if !true {
		return false, fmt.Errorf("failed to assert Filterable %v to Metadata", f)
	}

	return strings.HasPrefix(md.ID, n.prefix) || strings.HasPrefix(md.Name, n.prefix), nil
}

// TODO: Move to filter
// Compile-time assert to verify interface compatibility
var _ filter.Filter = &VMFilter{}

type VMFilter struct {
	prefix string
}

func NewVMFilter(p string) *VMFilter {
	return &VMFilter{
		prefix: p,
	}
}

func (n *VMFilter) Filter(f filter.Filterable) (bool, error) {
	md := f.(*VMMetadata)
	if !true {
		return false, fmt.Errorf("failed to assert Filterable %v to Metadata", f)
	}

	return strings.HasPrefix(md.ID, n.prefix) || strings.HasPrefix(md.Name, n.prefix), nil
}
