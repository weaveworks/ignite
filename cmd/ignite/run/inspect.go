package run

import (
	"bytes"
	"fmt"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/filter"

	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

type InspectFlags struct {
	YAMLOutput bool
}

type inspectOptions struct {
	*InspectFlags
	object meta.Object
}

func (i *InspectFlags) NewInspectOptions(k, objectMatch string) (*inspectOptions, error) {
	var err error
	var kind meta.Kind
	io := &inspectOptions{InspectFlags: i}

	switch k {
	case meta.KindImage.Lower():
		kind = meta.KindImage
	case meta.KindKernel.Lower():
		kind = meta.KindKernel
	case meta.KindVM.Lower():
		kind = meta.KindVM
	default:
		return nil, fmt.Errorf("unrecognized kind: %s", k)
	}

	if io.object, err = client.Dynamic(kind).Find(filter.NewIDNameFilter(objectMatch)); err != nil {
		return nil, err
	}

	return io, nil
}

func Inspect(io *inspectOptions) error {
	// Choose the encoder
	encodeFunc := scheme.Serializer.EncodeJSON
	if io.YAMLOutput {
		encodeFunc = scheme.Serializer.EncodeYAML
	}

	// Encode the object with the selected encoder
	b, err := encodeFunc(io.object)
	if err != nil {
		return err
	}

	// Print the encoded object
	fmt.Println(string(bytes.TrimSpace(b)))
	return nil
}
