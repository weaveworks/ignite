package main

import (
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/filter"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
)

type options struct {
	vm *vmmd.VM
}

func NewOptions(vmMatch string) (*options, error) {
	co := &options{}

	if vm, err := client.VMs().Find(filter.NewIDNameFilter(vmMatch)); err == nil {
		co.vm = &vmmd.VM{vm}
	} else {
		return nil, err
	}

	return co, nil
}
