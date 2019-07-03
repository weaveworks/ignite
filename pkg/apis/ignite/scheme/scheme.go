package scheme

import (
	"k8s.io/apimachinery/pkg/runtime"
	k8sserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/storage/serializer"
)

var (
	// Scheme is the runtime.Scheme to which all types are registered.
	Scheme = runtime.NewScheme()

	// codecs provides access to encoding and decoding for the scheme.
	// codecs is private, as Serializer will be used for all higher-level encoding/decoding
	codecs = k8sserializer.NewCodecFactory(Scheme)

	// Serializer provides high-level encoding/decoding functions
	Serializer = serializer.NewSerializer(Scheme, &codecs)
)

func init() {
	AddToScheme(Scheme)
}

// AddToScheme builds the scheme using all known versions of the api.
func AddToScheme(scheme *runtime.Scheme) {
	utilruntime.Must(v1alpha1.AddToScheme(Scheme))
	utilruntime.Must(scheme.SetVersionPriority(v1alpha1.SchemeGroupVersion))
}
