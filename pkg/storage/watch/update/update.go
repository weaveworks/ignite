package update

import (
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/storage"
)

// Update bundles an FileEvent with an
// APIType for Storage retrieval.
type Update struct {
	Event   ObjectEvent
	APIType meta.Object
}

// AssociatedUpdate bundles together an Update and a Storage
// implementation. This is used by SyncStorage to query the
// correct Storage for the updated Object.
type AssociatedUpdate struct {
	Update
	Storage storage.Storage
}
