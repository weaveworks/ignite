package run

import (
	"bytes"
	"fmt"
	"strings"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/filter"
	"github.com/weaveworks/ignite/pkg/providers"
)

type InspectFlags struct {
	OutputFormat string
}

type inspectOptions struct {
	*InspectFlags
	object meta.Object
}

func (i *InspectFlags) NewInspectOptions(k, objectMatch string) (*inspectOptions, error) {
	var err error
	var kind meta.Kind
	io := &inspectOptions{InspectFlags: i}

	switch strings.ToLower(k) {
	case api.KindImage.Lower():
		kind = api.KindImage
	case api.KindKernel.Lower():
		kind = api.KindKernel
	case api.KindVM.Lower():
		kind = api.KindVM
	default:
		return nil, fmt.Errorf("unrecognized kind: %q", k)
	}

	if io.object, err = providers.Client.Dynamic(kind).Find(filter.NewIDNameFilter(objectMatch)); err != nil {
		return nil, err
	}

	return io, nil
}

func Inspect(io *inspectOptions) error {
	var b []byte
	var err error

	// Select the encoder and encode the object with it
	switch io.OutputFormat {
	case "json":
		b, err = scheme.Serializer.EncodeJSON(io.object)
	case "yaml":
		b, err = scheme.Serializer.EncodeYAML(io.object)
	default:
		err = fmt.Errorf("unrecognized output format: %q", io.OutputFormat)
	}

	if err != nil {
		return err
	}

	// Print the encoded object
	fmt.Println(string(bytes.TrimSpace(b)))
	return nil
}
