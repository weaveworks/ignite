package metadata

import (
	"crypto/rand"
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/filter"
	"github.com/weaveworks/ignite/pkg/storage/filterer"
	"github.com/weaveworks/ignite/pkg/util"
)

var (
	nameRegex = regexp.MustCompile(`^[a-z-_0-9.:/@]*$`)
	uidRegex  = regexp.MustCompile(`^[a-z0-9]{16}$`)
)

type Metadata interface {
	meta.Object
	ObjectPath() string
	Save() error
}

// InitObject shall be run in all New{Kind} methods. This method should be
// private once possible (the runtime objects should move into this pkg)
func InitObject(obj meta.Object, c *client.Client) error {
	if obj == nil {
		return fmt.Errorf("object cannot be nil when initializing runtime data")
	}

	if c == nil {
		c = client.DefaultClient
	}

	// Default the object
	scheme.Scheme.Default(obj)

	// Generate or validate the given UID, if any
	if err := processUID(obj, c); err != nil {
		return err
	}

	// Generate or validate the given name, if any
	if err := processName(obj, c); err != nil {
		return err
	}

	return nil
}

// processUID a new 8-byte ID and handles directory creation/deletion
func processUID(obj meta.Object, c *client.Client) error {
	uid := obj.GetUID().String()

	// Validate the given UID if set
	if len(uid) > 0 {
		// Verify that if specified
		if !uidRegex.MatchString(uid) {
			return fmt.Errorf("invalid UID %q: does not match required format %s", uid, uidRegex.String())
		}

		// Make sure there isn't any duplicate names
		if err := verifyUIDOrName(c, uid, obj.GetKind()); err != nil {
			return err
		}
	} else {
		// No UID set, generate one
		var uidBytes []byte
		for {
			uidBytes = make([]byte, 8)
			if _, err := rand.Read(uidBytes); err != nil {
				return fmt.Errorf("failed to generate ID: %v", err)
			}

			// Convert the byte slice to a string literally
			uid = fmt.Sprintf("%x", uidBytes)

			// If the generated UID is unique break the generator loop
			if err := verifyUIDOrName(c, uid, obj.GetKind()); err == nil {
				// Set the generated UID to the object
				obj.SetUID(meta.UID(uid))
				break
			}
		}
	}

	// Create the directory for the specified UID
	// TODO: Move this kind of functionality into pkg/storage
	dir := path.Join(constants.DATA_DIR, obj.GetKind().Lower(), uid)
	if err := os.MkdirAll(dir, constants.DATA_DIR_PERM); err != nil {
		return fmt.Errorf("failed to create directory for ID %q: %v", uid, err)
	}

	return nil
}

func processName(obj meta.Object, c *client.Client) error {
	name := obj.GetName()
	kind := obj.GetKind()

	// Enforce a latest tag for images and kernels. Also,
	// images and kernels must have their name set at this stage
	if kind == api.KindImage || kind == api.KindKernel {
		if len(name) == 0 {
			// this should not happen, programmer error
			return fmt.Errorf("%s name must not be unset", kind.String())
		}
	} else if len(name) == 0 { // If some other kind's name is empty, set a random name
		name = util.RandomName()
	}

	// Validate the name with the regexp
	if !nameRegex.MatchString(name) {
		return fmt.Errorf("invalid name %q: does not match required format %s", name, nameRegex.String())
	}

	// Make sure there isn't any duplicate names
	if err := verifyUIDOrName(c, name, kind); err != nil {
		return err
	}

	// write the desired name to the object
	obj.SetName(name)
	return nil
}

func verifyUIDOrName(c *client.Client, match string, kind meta.Kind) error {
	_, err := c.Dynamic(kind).Find(filter.NewIDNameFilter(match))
	switch err.(type) {
	case *filterer.NonexistentError:
		// The id/name is unique, no error
		return nil
	case nil, *filterer.AmbiguousError:
		// The ambiguous error can only occur if someone manually created two Objects with the same name
		return fmt.Errorf("invalid %s id/name %q: already exists", kind, match)
	default:
		return err
	}
}
