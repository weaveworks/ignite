package metadata

import (
	"fmt"
	"log"
	"os"

	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

type Metadata interface {
	meta.Object
	TypePath() string
	ObjectPath() string
	Load() error
	Save() error
}

func Remove(md Metadata, quiet bool) error {
	if err := os.RemoveAll(md.ObjectPath()); err != nil {
		return fmt.Errorf("unable to remove directory for %s %q: %v", md.GetKind(), md.GetUID(), err)
	}

	if quiet {
		fmt.Println(md.GetUID())
	} else {
		log.Printf("Removed %s with name %q and ID %q", md.GetKind(), md.GetName(), md.GetUID())
	}

	return nil
}
