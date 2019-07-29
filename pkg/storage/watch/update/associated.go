package update

import (
	"github.com/weaveworks/ignite/pkg/storage"
)

// AssociatedUpdate bundles together an Update and a Storage
// implementation. This is used by SyncStorage to query the
// correct Storage for the updated Object.
type AssociatedUpdate struct {
	Update
	Storage storage.Storage
}
