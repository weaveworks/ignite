package storage

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

// This can be used to match either of the keys
type AnyKey fmt.Stringer

// KindKey represents the internal format of Kind virtual paths
type KindKey struct {
	meta.Kind
}

// Key represents the internal format of Object virtual paths
type Key struct {
	KindKey
	meta.UID
}

// NewKindKey generates a new virtual path Key for a Kind
func NewKindKey(kind meta.Kind) KindKey {
	return KindKey{
		kind,
	}
}

// NewKey generates a new virtual path Key for an Object
func NewKey(kind meta.Kind, uid meta.UID) Key {
	return Key{
		NewKindKey(kind),
		uid,
	}
}

// ParseKey parses the given string and returns a Key
func ParseKey(input string) (k Key, err error) {
	splitInput := strings.Split(filepath.Clean(input), string(os.PathSeparator))
	if len(splitInput) != 2 {
		err = fmt.Errorf("invalid input for key parsing: %s", input)
	} else {
		k.Kind = parseKind(splitInput[0])
		k.UID = meta.UID(splitInput[1])
	}

	return
}

// String returns the virtual path for the Kind
func (k KindKey) String() string {
	return k.Lower()
}

// String returns the virtual path for the Object
func (k Key) String() string {
	return path.Join(k.KindKey.String(), k.UID.String())
}

// ToKindKey creates a KindKey out of a Key
func (k Key) ToKindKey() KindKey {
	return k.KindKey
}

// TODO: Move to meta
func parseKind(input string) meta.Kind {
	b := bytes.ToUpper([]byte(input))

	// Leave TLAs as uppercase
	if len(b) > 3 {
		b = append(b[:1], bytes.ToLower(b[1:])...)
	}

	return meta.Kind(b)
}
