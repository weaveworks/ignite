package patch

import (
	"fmt"

	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

// The default serializer used here. In the future we maybe want to
// make this configurable
var serializer = scheme.Serializer

// Create is a helper that creates a patch out of the change made in applyFn
func Create(new meta.Object, applyFn func(meta.Object) error) ([]byte, error) {
	old := new.DeepCopyObject().(meta.Object)

	oldbytes, err := serializer.EncodeJSON(old)
	if err != nil {
		return nil, err
	}

	emptyobj, err := serializer.Scheme().New(old.GroupVersionKind())
	if err != nil {
		return nil, err
	}

	if err := applyFn(new); err != nil {
		return nil, err
	}

	newbytes, err := serializer.EncodeJSON(new)
	if err != nil {
		return nil, err
	}

	patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldbytes, newbytes, emptyobj)
	if err != nil {
		return nil, fmt.Errorf("CreateTwoWayMergePatch failed: %v", err)
	}

	return patchBytes, nil
}

func Apply(original, patch []byte, gvk schema.GroupVersionKind) ([]byte, error) {
	emptyobj, err := serializer.Scheme().New(gvk)
	if err != nil {
		return nil, err
	}

	return strategicpatch.StrategicMergePatch(original, patch, emptyobj)
}
