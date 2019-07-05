package metadata

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"path"

	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/logs"
	"github.com/weaveworks/ignite/pkg/util"
)

var success = make(map[Metadata]bool)

// Creates a new 8-byte ID and handles directory creation/deletion
func NewUID(md Metadata, input meta.UID) error {
	// If a valid ID is given, don't overwrite it
	md.SetUID(input)
	if input != "" {
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
		idPath = path.Join(md.TypePath(), id)
		if exists, _ := util.PathExists(idPath); !exists {
			break
		}
	}

	// Create the directory for the ID
	if err := os.MkdirAll(idPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory for ID %q: %v", id, err)
	}

	// Set the generated ID
	md.SetUID(meta.UID(id))
	return nil
}

// silent specifies if the ID should be printed, when chaining commands
// silence all but the last command to print the ID only once
func Cleanup(md Metadata, silent bool) error {
	// If success has not been confirmed, remove the generated directory
	if !success[md] {
		if !logs.Quiet {
			log.Printf("Removed %s with name %q and ID %q", md.GetKind(), md.GetName(), md.GetUID())
		} else if !silent {
			fmt.Println(md.GetUID())
		}
		return client.Dynamic(md.GetKind().String()).Delete(md.GetUID())
	}

	if !logs.Quiet {
		log.Printf("Created %s with ID %q and name %q", md.GetKind(), md.GetUID(), md.GetName())
	} else if !silent {
		fmt.Println(md.GetUID())
	}

	return nil
}

// Should be returned as the last command when creating objects
func Success(md Metadata) error {
	success[md] = true
	return nil
}
