package update

import (
	"github.com/weaveworks/ignite/pkg/storage"
)

type AssociatedUpdate struct {
	Update
	Storage storage.Storage
}
