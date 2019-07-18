package providers

import "github.com/weaveworks/ignite/pkg/client"

// Client is the default client that can be easily used
var Client *client.Client

func SetClient() error {
	Client = client.NewClient(Storage)
	return nil
}
