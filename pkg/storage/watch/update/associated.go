package update

import (
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/storage"
)

type AssociatedUpdate struct {
	Event   Event
	APIType meta.Object
	Storage storage.Storage
}
