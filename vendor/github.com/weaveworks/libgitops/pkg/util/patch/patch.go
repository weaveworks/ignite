package patch

import (
	"fmt"
	"io/ioutil"

	"github.com/weaveworks/libgitops/pkg/runtime"
	"github.com/weaveworks/libgitops/pkg/serializer"
	kruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

type Patcher interface {
	Create(new runtime.Object, applyFn func(runtime.Object) error) ([]byte, error)
	Apply(original, patch []byte, gvk schema.GroupVersionKind) ([]byte, error)
	ApplyOnFile(filePath string, patch []byte, gvk schema.GroupVersionKind) error
}

func NewPatcher(s serializer.Serializer) Patcher {
	return &patcher{serializer: s}
}

type patcher struct {
	serializer serializer.Serializer
}

// Create is a helper that creates a patch out of the change made in applyFn
func (p *patcher) Create(new runtime.Object, applyFn func(runtime.Object) error) ([]byte, error) {
	old := new.DeepCopyObject().(runtime.Object)

	oldbytes, err := p.serializer.EncodeJSON(old)
	if err != nil {
		return nil, err
	}

	emptyobj, err := p.serializer.Scheme().New(old.GroupVersionKind())
	if err != nil {
		return nil, err
	}

	if err := applyFn(new); err != nil {
		return nil, err
	}

	newbytes, err := p.serializer.EncodeJSON(new)
	if err != nil {
		return nil, err
	}

	patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldbytes, newbytes, emptyobj)
	if err != nil {
		return nil, fmt.Errorf("CreateTwoWayMergePatch failed: %v", err)
	}

	return patchBytes, nil
}

func (p *patcher) Apply(original, patch []byte, gvk schema.GroupVersionKind) ([]byte, error) {
	emptyobj, err := p.serializer.Scheme().New(gvk)
	if err != nil {
		return nil, err
	}

	b, err := strategicpatch.StrategicMergePatch(original, patch, emptyobj)
	if err != nil {
		return nil, err
	}

	return p.serializerEncode(b)
}

func (p *patcher) ApplyOnFile(filePath string, patch []byte, gvk schema.GroupVersionKind) error {
	oldContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	newContent, err := p.Apply(oldContent, patch, gvk)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, newContent, 0644)
}

// StrategicMergePatch returns an unindented, unorganized JSON byte slice,
// this helper takes that as an input and returns the same JSON re-encoded
// with the serializer so it conforms to a runtime.Object
// TODO: Just use encoding/json.Indent here instead?
func (p *patcher) serializerEncode(input []byte) (result []byte, err error) {
	var obj kruntime.Object
	if obj, err = p.serializer.Decode(input, true); err == nil {
		result, err = p.serializer.EncodeJSON(obj)
	}

	return
}
