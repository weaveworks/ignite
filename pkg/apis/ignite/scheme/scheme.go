package scheme

import (
	"io/ioutil"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
)

var (
	// Scheme is the runtime.Scheme to which all types are registered.
	Scheme = runtime.NewScheme()

	// Codecs provides access to encoding and decoding for the scheme.
	Codecs = serializer.NewCodecFactory(Scheme)
)

func init() {
	AddToScheme(Scheme)
}

// AddToScheme builds the scheme using all known versions of the api.
func AddToScheme(scheme *runtime.Scheme) {
	utilruntime.Must(v1alpha1.AddToScheme(Scheme))
	utilruntime.Must(scheme.SetVersionPriority(v1alpha1.SchemeGroupVersion))
}

// DecodeFileInto takes a file path and a target object to serialize the data into
func DecodeFileInto(filePath string, obj runtime.Object) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return DecodeInto(content, obj)
}

// DecodeInto takes byte content and a target object to serialize the data into
func DecodeInto(content []byte, obj runtime.Object) error {
	return runtime.DecodeInto(Codecs.UniversalDecoder(), content, obj)
}

// EncodeYAML encodes the specified object for a specific version to YAML bytes
func EncodeYAML(obj runtime.Object, groupVersion schema.GroupVersion) ([]byte, error) {
	serializerInfo, _ := runtime.SerializerInfoForMediaType(Codecs.SupportedMediaTypes(), runtime.ContentTypeYAML)
	return runtime.Encode(Codecs.EncoderForVersion(serializerInfo.Serializer, groupVersion), obj)
}

// EncodeYAML encodes the specified object for a specific version to pretty JSON bytes
func EncodeJSON(obj runtime.Object, groupVersion schema.GroupVersion) ([]byte, error) {
	serializerInfo, _ := runtime.SerializerInfoForMediaType(Codecs.SupportedMediaTypes(), runtime.ContentTypeJSON)
	return runtime.Encode(Codecs.EncoderForVersion(serializerInfo.PrettySerializer, groupVersion), obj)
}
