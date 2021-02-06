package run

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/filter"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/util"
)

// PsFlags contains the flags supported by ps.
type PsFlags struct {
	All            bool
	Filter         string
	TemplateFormat string
}

type PsOptions struct {
	*PsFlags
	allVMs []*api.VM
}

// NewPsOptions constructs and returns PsOptions.
func (pf *PsFlags) NewPsOptions() (po *PsOptions, err error) {
	po = &PsOptions{PsFlags: pf}
	po.allVMs, err = providers.Client.VMs().FindAll(filter.NewVMFilterAll("", po.All))
	// If the storage is uninitialized, avoid failure and continue with empty
	// VM list.
	if err != nil && strings.Contains(err.Error(), "no such file or directory") {
		err = nil
	}
	return
}

// Ps filters and renders the VMs based on the PsOptions.
func Ps(po *PsOptions) error {
	var filters *filter.MultipleMetaFilter
	var err error
	var filtering bool
	if po.PsFlags.Filter != "" {
		filtering = true
		filters, err = filter.GenerateMultipleMetadataFiltering(po.PsFlags.Filter)
		if err != nil {
			return err
		}
	}

	filteredVMs := []*api.VM{}

	for _, vm := range po.allVMs {
		isExpectedVM := true
		if filtering {
			isExpectedVM, err = filters.AreExpected(vm)
			if err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
		if isExpectedVM {
			filteredVMs = append(filteredVMs, vm)
		}
	}

	// If template format is specified, render the template.
	if po.PsFlags.TemplateFormat != "" {
		// Parse the template format.
		tmpl, err := template.New("").Parse(po.PsFlags.TemplateFormat)
		if err != nil {
			return fmt.Errorf("failed to parse template: %v", err)
		}

		// Render the template with the filtered VMs.
		for _, vm := range filteredVMs {
			o := &bytes.Buffer{}
			if err := tmpl.Execute(o, vm); err != nil {
				return fmt.Errorf("failed rendering template: %v", err)
			}
			fmt.Println(o.String())
		}
		return nil
	}

	o := util.NewOutput()
	defer o.Flush()

	o.Write("VM ID", "IMAGE", "KERNEL", "SIZE", "CPUS", "MEMORY", "CREATED", "STATUS", "IPS", "PORTS", "NAME")
	for _, vm := range filteredVMs {
		o.Write(vm.GetUID(), vm.Spec.Image.OCI, vm.Spec.Kernel.OCI,
			vm.Spec.DiskSize, vm.Spec.CPUs, vm.Spec.Memory, formatCreated(vm), formatStatus(vm), vm.Status.Network.IPAddresses,
			vm.Spec.Network.Ports, vm.GetName())
	}

	return nil
}

func formatCreated(vm *api.VM) string {
	created := vm.GetCreated()

	var suffix string
	if !created.IsZero() {
		suffix = " ago"
	}

	return fmt.Sprint(created, suffix)
}

func formatStatus(vm *api.VM) string {
	if vm.Running() {
		return fmt.Sprintf("Up %s", vm.Status.StartTime)
	}

	return "Stopped"
}
