package metadata

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/weaveworks/ignite/pkg/source"
	"github.com/weaveworks/ignite/pkg/logs"
	"github.com/weaveworks/ignite/pkg/util"
)

type ID struct {
	string
	success bool
}

// Compile-time asserts to verify interface compatibility
var _ fmt.Stringer = &ID{}
var _ json.Marshaler = &ID{}
var _ json.Unmarshaler = &ID{}

func (id *ID) String() string {
	return id.string
}

func (id *ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.string)
}

func (id *ID) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	*id = ID{string: s}

	return nil
}

func (id *ID) Equal(other *ID) bool {
	return id.string == other.string
}

func IDFromSource(src source.Source) *ID {
	return &ID{
		string: src.ID(),
	}
}

// Creates a new 8-byte ID and handles directory creation/deletion
func (md *Metadata) newID() error {
	// If there's already an ID set, don't overwrite it
	if md.ID != nil {
		return nil
	}

	var id string
	var idPath string
	var idBytes []byte

	for {
		idBytes = make([]byte, 8)
		if _, err := rand.Read(idBytes); err != nil {
			return fmt.Errorf("failed to generate ID: %v", err)
		}

		// Convert the byte slice to a string literally
		id = fmt.Sprintf("%x", idBytes)

		// If the generated ID is unique break the generator loop
		idPath = path.Join(md.Type.Path(), id)
		if exists, _ := util.PathExists(idPath); !exists {
			break
		}
	}

	// Create the directory for the ID
	if err := os.MkdirAll(idPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory for ID %q: %v", id, err)
	}

	// Set the generated ID
	md.ID = &ID{string: id}
	return nil
}

// silent specifies if the ID should be printed, when chaining commands
// silence all but the last command to print the ID only once
func (md *Metadata) Cleanup(silent bool) error {
	// If success has not been confirmed, remove the generated directory
	if !md.ID.success {
		return md.Remove(logs.Quiet)
	}

	if !logs.Quiet {
		log.Printf("Created %s with ID %q and name %q", md.Type, md.ID, md.Name)
	} else if !silent {
		fmt.Println(md.ID)
	}

	return nil
}

// Should be returned as the last command when creating objects
func (md *Metadata) Success() error {
	md.ID.success = true
	return nil
}
