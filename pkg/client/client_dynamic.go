package client

import (
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/storage"
	"github.com/weaveworks/ignite/pkg/storage/filterer"
)

// DynamicClient is an interface for accessing API types generically
type DynamicClient interface {
	// Get returns an Object matching the UID from the storage
	Get(meta.UID) (meta.Object, error)
	// Set saves an Object into the persistent storage
	Set(meta.Object) error
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
func (c *Client) Dynamic(kind meta.Kind) (dc DynamicClient) {
	var ok bool
	if dc, ok = c.dynamicClients[kind]; !ok {
		dc = newDynamicClient(c.storage, kind)
		c.dynamicClients[kind] = dc
	}

	return
}

// Dynamic is a shorthand for accessing the DynamicClient using the default client
func Dynamic(kind meta.Kind) DynamicClient {
	return DefaultClient.Dynamic(kind)
}

// dynamicClient is a struct implementing the DynamicClient interface
// It uses a shared storage instance passed from the Client together with its own Filterer
type dynamicClient struct {
	storage  storage.Storage
	kind     meta.Kind
	filterer *filterer.Filterer
}

// newDynamicClient builds the dynamicClient struct using the storage implementation and a new Filterer
func newDynamicClient(s storage.Storage, kind meta.Kind) DynamicClient {
	return &dynamicClient{
		storage:  s,
		kind:     kind,
		filterer: filterer.NewFilterer(s),
	}
}

// Get returns an Object based the given UID
func (c *dynamicClient) Get(uid meta.UID) (meta.Object, error) {
	return c.storage.GetByID(c.kind, uid)
}

// Set saves an Object into the persistent storage
func (c *dynamicClient) Set(resource meta.Object) error {
	return c.storage.Set(resource)
}

// Find returns an Object based on a given Filter
func (c *dynamicClient) Find(filter filterer.BaseFilter) (meta.Object, error) {
	return c.filterer.Find(c.kind, filter)
}

// FindAll returns multiple Objects based on a given Filter
func (c *dynamicClient) FindAll(filter filterer.BaseFilter) ([]meta.Object, error) {
	return c.filterer.FindAll(c.kind, filter)
}

// Delete deletes the Object from the storage
func (c *dynamicClient) Delete(uid meta.UID) error {
	return c.storage.Delete(c.kind, uid)
}

// List returns a list of all Objects available
func (c *dynamicClient) List() ([]meta.Object, error) {
	return c.storage.List(c.kind)
}
