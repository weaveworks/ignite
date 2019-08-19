package patch

import (
	"fmt"
	"io/ioutil"

	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/serializer"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

type Patcher interface {
	Create(new meta.Object, applyFn func(meta.Object) error) ([]byte, error)
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
func (p *patcher) Create(new meta.Object, applyFn func(meta.Object) error) ([]byte, error) {
	old := new.DeepCopyObject().(meta.Object)

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

	return strategicpatch.StrategicMergePatch(original, patch, emptyobj)
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
