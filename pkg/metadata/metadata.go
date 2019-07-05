package metadata

import (
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

type Metadata interface {
	meta.Object
	TypePath() string
	ObjectPath() string
	Load() error
	Save() error
}
