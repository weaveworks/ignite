package metadata

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/util"
	"io/ioutil"
	"os"
	"path"
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
	VM:     "vm",
}

func (x ObjectType) MarshalJSON() ([]byte, error) {
	return json.Marshal(ObjectTypeLookup[x])
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

type ObjectData interface{}

type Metadata struct {
	ID         string     `json:"ID"`
	Name       string     `json:"Name"`
	Type       ObjectType `json:"Type"`
	ObjectData `json:"ObjectData"`
}

func (md *Metadata) ObjectPath() string {
	var prefix string

	switch md.Type {
	case Image:
		prefix = constants.IMAGE_DIR
	case Kernel:
		prefix = constants.KERNEL_DIR
	case VM:
		prefix = constants.VM_DIR
	}

	return path.Join(prefix, md.ID)
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
	t := ObjectTypeLookup[md.Type]

	if !util.DirExists(p) {
		return fmt.Errorf("nonexistent %q object: %s", t, md.ID)
	}

	f := path.Join(p, constants.METADATA)

	if !util.FileExists(f) {
		return fmt.Errorf("metadata file missing for %q object: %s", t, md.ID)
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
