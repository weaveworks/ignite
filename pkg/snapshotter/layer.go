package snapshotter

import "github.com/weaveworks/ignite/pkg/dm"

type layer struct {
	device *dm.Device
}

func newLayer(device *dm.Device) layer {
	return layer{
		device: device,
	}
}
