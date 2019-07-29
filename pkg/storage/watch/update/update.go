package update

import (
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

// Update bundles an Event with an
// APIType for Storage retrieval.
type Update struct {
	Event   Event
	APIType meta.Object
}
