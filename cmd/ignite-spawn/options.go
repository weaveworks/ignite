package main

import (
	"github.com/weaveworks/ignite/pkg/metadata/loader"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
)

type options struct {
	vm *vmmd.VMMetadata
}

func NewOptions(l *loader.ResLoader, vmMatch string) (*options, error) {
	co := &options{}

	if allVMS, err := l.VMs(); err == nil {
		if co.vm, err = allVMS.MatchSingle(vmMatch); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return co, nil
}
