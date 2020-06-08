package update

import (
	"github.com/weaveworks/libgitops/pkg/runtime"
	"github.com/weaveworks/libgitops/pkg/storage"
)

// Update bundles an FileEvent with an
// APIType for Storage retrieval.
type Update struct {
	Event   ObjectEvent
	APIType runtime.Object
}

// AssociatedUpdate bundles together an Update and a Storage
// implementation. This is used by SyncStorage to query the
// correct Storage for the updated Object.
type AssociatedUpdate struct {
	Update
	Storage storage.Storage
}
