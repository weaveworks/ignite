package update

import (
	"github.com/weaveworks/gitops-toolkit/pkg/runtime"
	"github.com/weaveworks/gitops-toolkit/pkg/storage"
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
