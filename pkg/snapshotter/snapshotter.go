package snapshotter

import (
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/util"
	"path"

	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	ignitemeta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/dm"
)

// A snapshotter Object represents the internal storage format of a device
// that binds the device mapper device to the object metadata
// TODO: The object needs to be an interface, otherwise the conversions get out of control
//type Object struct {
//	device *dm.Device
//	object ignitemeta.Object
//	parent *Object
//}

type Object interface {
	device() *dm.Device
	object() ignitemeta.Object
	parent() Object
}

type Layer interface {
}

// The metadata is behind a private field, which enables this getter to load it on-demand
func (o *Object) GetMetaObject() (ignitemeta.Object, error) {
	if o.object == nil && len(o.device.MetadataPath) > 0 {
		// TODO: !!!!!!!! This won't work! We're decoding into a generic interface!
		if err := scheme.DecodeFileInto(o.device.MetadataPath, o.object); err != nil {
			return nil, err
		}
	}

	return o.object, nil
}

// This is for filtering functions to resolve parents of images
func (o *Object) ChildOf(i *Image) bool {
	if o.device == i.device {
		return true
	}

	if o.parent != nil {
		return o.parent.ChildOf(i)
	}

	return false
}

// Snapshotter abstracts the device mapper pool and provides convenience methods
// It's also responsible for (de)serializing the pool
type Snapshotter struct {
	pool   *dm.Pool
	layers []Layer
	//
}

// NewSnapshotter creates a new Snapshotter with a new Pool
// or loads an existing configuration if it exists
// TODO: No support for physical backing devices for now
func NewSnapshotter() (*Snapshotter, error) {
	p := path.Join(constants.SNAPSHOTTER_DIR, constants.METADATA)
	s := &Snapshotter{}

	// If the metadata doesn't exist, return a new Snapshotter
	if !util.FileExists(p) {
		pool := &v1alpha1.Pool{}
		v1alpha1.SetObjectDefaults_Pool(pool)
		s.pool = dm.NewPool(pool)
		return s, nil
	}

	// Load the pool configuration
	if err := scheme.DecodeFileInto(p, s.pool); err != nil {
		return nil, err
	}

	// Allocate Objects for each device
	s.objects = make([]Object, s.pool.Size())
	for i := 0; i < len(s.objects); i++ {
		s.objects[i] = nil
	}

	// Create Objects from each device
	_ = s.pool.ForDevices(func(id ignitemeta.DMID, device *dm.Device) error {
		s.objects[id.Index()].device() = device
		s.objects[id.Index()].parent() = s.objects[device.Parent.Index()]

		return nil
	})

	return s, nil
}

// createObject creates a whole new object
func (s *Snapshotter) createObject() *Object {
	device, err := i.device.CreateSnapshot(size, "")
	if err != nil {
		return nil, err
	}
}
