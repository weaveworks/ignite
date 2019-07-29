package update

import (
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

type Update struct {
	Event   Event
	APIType meta.Object
}
