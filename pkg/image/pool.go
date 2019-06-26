package image

import (
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/dm"
)

// We need a pool subset which loads the given UID and it's children
// This is for per-image access to the same pool

type ImagePool struct {
	pool    v1alpha1.PoolData
	devices []dm.Device
}
