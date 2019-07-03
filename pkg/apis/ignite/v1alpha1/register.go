package api

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// SchemeBuilder the schema builder
	SchemeBuilder = runtime.NewSchemeBuilder(
		addKnownTypes,
		addDefaultingFuncs,
	)

	localSchemeBuilder = &SchemeBuilder
	AddToScheme        = localSchemeBuilder.AddToScheme
)

const (
	// GroupName is the group name use in this package
	GroupName = "ignite.weave.works"

	// VMKind returns the kind for the VM API type
	VMKind = "VM"
	// KernelKind returns the kind for the Kernel API type
	KernelKind = "Kernel"
	// PoolKind returns the kind for the Pool API type
	PoolKind = "Pool"
	// ImageKind returns the kind for the Image API type
	ImageKind = "Image"
)

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{
	Group:   GroupName,
	Version: "api",
}
var internalGV = schema.GroupVersion{
	Group:   GroupName,
	Version: runtime.APIVersionInternal,
}

// Adds the list of known types to the given scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&VM{},
		&Kernel{},
		&Pool{},
		&Image{},
	)
	// TODO: This is a hack, but for now it's sufficient.
	// Eventually, we should break this out to a real internal API package
	scheme.AddKnownTypes(internalGV,
		&VM{},
		&Kernel{},
		&Pool{},
		&Image{},
	)

	return nil
}
