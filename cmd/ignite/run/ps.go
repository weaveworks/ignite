package run

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/filter"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/runtime"
	containerdruntime "github.com/weaveworks/ignite/pkg/runtime/containerd"
	dockerruntime "github.com/weaveworks/ignite/pkg/runtime/docker"
	"github.com/weaveworks/ignite/pkg/util"
)

// runtimeRunningStatus is the status returned from the container runtimes when
// the VM container is in running state.
const runtimeRunningStatus = "running"

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

	outdatedVMs, err := fetchLatestStatus(filteredVMs)
	if err != nil {
		return err
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
			vm.Spec.DiskSize, vm.Spec.CPUs, vm.Spec.Memory, formatCreated(vm), formatStatus(vm, outdatedVMs), vm.Status.Network.IPAddresses,
			vm.Spec.Network.Ports, vm.GetName())
	}

	// Add a note at the bottom about the old manifest indicator in the status.
	if len(outdatedVMs) > 0 {
		o.Write("\nNOTE: The symbol * on the VM status indicates that the VM manifest on disk is not up-to-date with the actual VM status.")
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

func formatStatus(vm *api.VM, outdatedVMs map[string]bool) string {
	oldManifestIndicator := ""
	if _, ok := outdatedVMs[vm.Name]; ok {
		oldManifestIndicator = "*"
	}
	if vm.Running() {
		return fmt.Sprintf("%sUp %s", oldManifestIndicator, vm.Status.StartTime)
	}

	return oldManifestIndicator + "Stopped"
}

// fetchLatestStatus fetches the current status the VMs, updates the VM status
// in memory and returns a list of outdated VMs.
func fetchLatestStatus(vms []*api.VM) (outdatedVMs map[string]bool, err error) {
	outdatedVMs = map[string]bool{}

	// Container runtime clients. These clients are lazy initialized based on
	// the VM's runtime.
	var containerdClient, dockerClient runtime.Interface

	// Iterate through the VMs, fetching the actual status from the runtime.
	for _, vm := range vms {
		// Skip VMs with no runtime info.
		if vm.Status.Runtime == nil {
			continue
		}
		containerID := vm.Status.Runtime.ID
		currentRunning := false

		// Runtime client of the VM.
		var vmRuntime runtime.Interface

		// Set the appropriate runtime client based on the VM runtime info.
		switch vm.Status.Runtime.Name {
		case runtime.RuntimeContainerd:
			if containerdClient == nil {
				containerdClient, err = containerdruntime.GetContainerdClient()
				if err != nil {
					return
				}
			}
			vmRuntime = containerdClient
		case runtime.RuntimeDocker:
			if dockerClient == nil {
				dockerClient, err = dockerruntime.GetDockerClient()
				if err != nil {
					return
				}
			}
			vmRuntime = dockerClient
		}

		// Inspect the VM container using the runtime client.
		ir, inspectErr := vmRuntime.InspectContainer(containerID)
		if inspectErr != nil {
			err = errors.Wrapf(inspectErr, "failed to inspect container for VM %s", containerID)
			return
		}

		// Set current running based on the container status result.
		if ir.Status == runtimeRunningStatus {
			currentRunning = true
		}

		// If current running status and the VM object status don't match, mark
		// it as an outdated VM and update the VM object staus in memory.
		// NOTE: Avoid updating the VM manifest on disk here. That'll be
		// indicated in the ps output.
		if currentRunning != vm.Status.Running {
			vm.Status.Running = currentRunning
			outdatedVMs[vm.Name] = true
		}
	}

	return
}
