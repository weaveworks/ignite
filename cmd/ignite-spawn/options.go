package main

import (
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/filter"
)

type options struct {
	vm *api.VM
}

func NewOptions(vmMatch string) (*options, error) {
	co := &options{}

	if vm, err := client.VMs().Find(filter.NewIDNameFilter(vmMatch)); err == nil {
		co.vm = vm
	} else {
		return nil, err
	}

	return co, nil
}
