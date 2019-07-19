package main

import (
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/filter"
	"github.com/weaveworks/ignite/pkg/providers"
)

type options struct {
	vm *api.VM
}

func NewOptions(vmMatch string) (*options, error) {
	co := &options{}

	if vm, err := providers.Client.VMs().Find(filter.NewIDNameFilter(vmMatch)); err == nil {
		co.vm = vm
	} else {
		return nil, err
	}

	return co, nil
}
