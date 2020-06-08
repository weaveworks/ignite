package storage

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/weaveworks/libgitops/pkg/runtime"
)

// This can be used to match either of the keys
type AnyKey fmt.Stringer

// KindKey represents the internal format of Kind virtual paths
type KindKey struct {
	runtime.Kind
}

// Key represents the internal format of Object virtual paths
type Key struct {
	KindKey
	runtime.UID
}

// NewKindKey generates a new virtual path Key for a Kind
func NewKindKey(kind runtime.Kind) KindKey {
	return KindKey{
		kind,
	}
}

// NewKey generates a new virtual path Key for an Object
func NewKey(kind runtime.Kind, uid runtime.UID) Key {
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
		k.Kind = runtime.ParseKind(splitInput[0])
		k.UID = runtime.UID(splitInput[1])
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
