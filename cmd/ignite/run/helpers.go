package run

import (
	"github.com/weaveworks/ignite/pkg/filter"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/providers"
)

// TODO: This
func getVMForMatch(vmMatch string) (*vmmd.VM, error) {
	apiVM, err := providers.Client.VMs().Find(filter.NewIDNameFilter(vmMatch))
	if err != nil {
		return nil, err
	}
	return vmmd.WrapVM(apiVM), nil
}

// TODO: This
func getVMsForMatches(vmMatches []string) ([]*vmmd.VM, error) {
	allVMs := make([]*vmmd.VM, 0, len(vmMatches))
	for _, match := range vmMatches {
		runVM, err := getVMForMatch(match)
		if err != nil {
			return nil, err
		}
		allVMs = append(allVMs, runVM)
	}
	return allVMs, nil
}

func getAllVMs() (allVMs []*vmmd.VM, err error) {
	allAPIVMs, err := providers.Client.VMs().FindAll(filter.NewAllFilter())
	if err != nil {
		return
	}
	allVMs = make([]*vmmd.VM, 0, len(allAPIVMs))
	for _, apiVM := range allAPIVMs {
		allVMs = append(allVMs, vmmd.WrapVM(apiVM))
	}
	return
}
