/*
Package client is a Go client for Ignite resources.

For example, to list running VMs (the equivalent of "ignite vm ls"), and update the
VM with a new IP address:

	package main
	import (
		"context"
		"fmt"
		"net"
		"github.com/weaveworks/ignite/pkg/client"
	)
	func main() {
		// List VMs managed by Ignite
		vmList, err := client.VMs().List()
		if err != nil {
			panic(err)
		}
		for _, vm := range vmList {
			// Modify the object
			vm.Status.IPAddresses = append(vm.Status.IPAddresses, net.IP{127,0,0,1})
			// Save the new VM state
			if err := client.VMs().Set(vm); err != nil {
				panic(err)
			}
		}

		// Get a specific image, and print its size
		myImage, err := client.Images().Get("my-image")
		if err != nil {
			panic(err)
		}
		fmt.Printf("Image my-vm size: %s\n", myImage.Spec.Source.Size.String())

		// Use the dynamic client
		myVM, err := client.Dynamic("VM").Get("my-vm")
		if err != nil {
			panic(err)
		}
		fmt.Printf("VM my-vm cpus: %d\n", myVM.Spec.CPUs)
	}

*/
package client

import (
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/storage"
)

// NewClient creates a client for the specified storage
func NewClient(s storage.Storage) *Client {
	return &Client{
		storage:        s,
		dynamicClients: map[meta.Kind]DynamicClient{},
	}
}

// Client is a struct providing high-level access to objects in a storage
// The resource-specific client interfaces are automatically generated based
// off client_resource_template.go. The auto-generation can be done with hack/client.sh
type Client struct {
	storage        storage.Storage
	vmClient       VMClient
	kernelClient   KernelClient
	imageClient    ImageClient
	dynamicClients map[meta.Kind]DynamicClient
}
