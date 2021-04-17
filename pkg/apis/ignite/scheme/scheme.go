package scheme

import (
	"k8s.io/apimachinery/pkg/runtime"
	k8sserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	"github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha2"
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha3"
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha4"
	"github.com/weaveworks/libgitops/pkg/serializer"
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
	utilruntime.Must(ignite.AddToScheme(Scheme))
	utilruntime.Must(v1alpha2.AddToScheme(Scheme))
	utilruntime.Must(v1alpha3.AddToScheme(Scheme))
	utilruntime.Must(v1alpha4.AddToScheme(Scheme))
	utilruntime.Must(scheme.SetVersionPriority(v1alpha4.SchemeGroupVersion))
}
