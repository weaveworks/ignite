package client

import (
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/storage"
)

// DynamicClient is an interface for accessing API types generically
type DynamicClient interface {
	storage.Cache

	// Get returns a Resource object based on a reference string; which can either
	// match the Resource's Name or UID, or be a prefix of the UID
	Get(ref string) (meta.Object, error)
	// Set saves a Resource into the persistent storage
	Set(meta.Object) error
	// Delete deletes the API object from the storage
	Delete(uid meta.UID) error
	// List returns a list of all Resources available
	List() ([]meta.Object, error)
}

// Dynamic returns the DynamicClient for the Client instance, for the specific kind
func (c *Client) Dynamic(kind string) DynamicClient {
	dc, ok := c.dynamicClients[kind]
	if !ok {
		c.dynamicClients[kind] = newDynamicClient(c.storage, kind)
		dc = c.dynamicClients[kind]
	}
	return dc
}

// Dynamic is a shorthand for accessing the DynamicClient using the default client
func Dynamic(kind string) DynamicClient {
	return DefaultClient.Dynamic(kind)
}

// dynamicClient is a struct implementing the DynamicClient interface
// It uses a shared storage instance passed from the Client
type dynamicClient struct {
	storage.Cache
	storage storage.Storage
	kind    string
}

// newDynamicClient builds the dynamicClient struct using the storage implementation
// It automatically fetches all metadata for all API types of the specific kind into the cache
func newDynamicClient(s storage.Storage, kind string) DynamicClient {
	c, err := s.GetCache(kind)
	if err != nil {
		panic(err)
	}
	return &dynamicClient{storage: s, Cache: c, kind: kind}
}

// Get returns a Resource object based on a reference string; which can either
// match the Resource's Name or UID, or be a prefix of the UID
func (c *dynamicClient) Get(ref string) (meta.Object, error) {
	meta, err := c.MatchOne(ref)
	if err != nil {
		return nil, err
	}
	obj, err := c.storage.GetByID(meta.Kind, meta.UID)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// Set saves a Resource into the persistent storage
func (c *dynamicClient) Set(resource meta.Object) error {
	return c.storage.Set(resource)
}

// Delete deletes the API object from the storage
func (c *dynamicClient) Delete(uid meta.UID) error {
	return c.storage.Delete(c.kind, uid)
}

// List returns a list of all Resources available
func (c *dynamicClient) List() ([]meta.Object, error) {
	return c.storage.List(c.kind)
}
