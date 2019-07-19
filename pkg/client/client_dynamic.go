package client

import (
	"fmt"

	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/storage"
	"github.com/weaveworks/ignite/pkg/storage/filterer"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

// DynamicClient is an interface for accessing API types generically
type DynamicClient interface {
	// New returns a new Object of its kind
	New() meta.Object
	// Get returns an Object matching the UID from the storage
	Get(meta.UID) (meta.Object, error)
	// Set saves an Object into the persistent storage
	Set(meta.Object) error
	// Patch performs a strategic merge patch on the object with
	// the given UID, using the byte-encoded patch given
	Patch(meta.UID, []byte) error
	// Find returns an Object based on the given filter, filters can
	// match e.g. the Object's Name, UID or a specific property
	Find(filter filterer.BaseFilter) (meta.Object, error)
	// FindAll returns multiple Objects based on the given filter, filters can
	// match e.g. the Object's Name, UID or a specific property
	FindAll(filter filterer.BaseFilter) ([]meta.Object, error)
	// Delete deletes an Object from the storage
	Delete(uid meta.UID) error
	// List returns a list of all Objects available
	List() ([]meta.Object, error)
}

// Dynamic returns the DynamicClient for the Client instance, for the specific kind
func (c *IgniteInternalClient) Dynamic(kind meta.Kind) (dc DynamicClient) {
	var ok bool
	gvk := c.gv.WithKind(kind.Title())
	if dc, ok = c.dynamicClients[gvk]; !ok {
		dc = newDynamicClient(c.storage, gvk)
		c.dynamicClients[gvk] = dc
	}

	return
}

// dynamicClient is a struct implementing the DynamicClient interface
// It uses a shared storage instance passed from the Client together with its own Filterer
type dynamicClient struct {
	storage  storage.Storage
	gvk      schema.GroupVersionKind
	filterer *filterer.Filterer
}

// newDynamicClient builds the dynamicClient struct using the storage implementation and a new Filterer
func newDynamicClient(s storage.Storage, gvk schema.GroupVersionKind) DynamicClient {
	return &dynamicClient{
		storage:  s,
		gvk:      gvk,
		filterer: filterer.NewFilterer(s),
	}
}

// New returns a new Object of its kind
func (c *dynamicClient) New() meta.Object {
	obj, err := c.storage.New(c.gvk)
	if err != nil {
		panic(fmt.Sprintf("Client.New must not return an error: %v", err))
	}
	return obj
}

// Get returns an Object based the given UID
func (c *dynamicClient) Get(uid meta.UID) (meta.Object, error) {
	return c.storage.Get(c.gvk, uid)
}

// Set saves an Object into the persistent storage
func (c *dynamicClient) Set(resource meta.Object) error {
	return c.storage.Set(c.gvk, resource)
}

// Patch performs a strategic merge patch on the object with
// the given UID, using the byte-encoded patch given
func (c *dynamicClient) Patch(uid meta.UID, patch []byte) error {
	return c.storage.Patch(c.gvk, uid, patch)
}

// Find returns an Object based on a given Filter
func (c *dynamicClient) Find(filter filterer.BaseFilter) (meta.Object, error) {
	return c.filterer.Find(c.gvk, filter)
}

// FindAll returns multiple Objects based on a given Filter
func (c *dynamicClient) FindAll(filter filterer.BaseFilter) ([]meta.Object, error) {
	return c.filterer.FindAll(c.gvk, filter)
}

// Delete deletes the Object from the storage
func (c *dynamicClient) Delete(uid meta.UID) error {
	return c.storage.Delete(c.gvk, uid)
}

// List returns a list of all Objects available
func (c *dynamicClient) List() ([]meta.Object, error) {
	return c.storage.List(c.gvk)
}
