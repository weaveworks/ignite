package run

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
const oldManifestIndicator = "*"

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
	if err != nil && os.IsNotExist(err) {
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

	endWarnings := []error{}
	outdatedVMs, errList := fetchLatestStatus(filteredVMs)
	if len(outdatedVMs) > 0 {
		endWarnings = append(
			endWarnings,
			fmt.Errorf("The symbol %s on the VM status indicates that the VM manifest on disk may not be up-to-date with the actual VM status from the container runtime", oldManifestIndicator),
		)
	}
	if len(errList) > 0 {
		endWarnings = append(endWarnings, errList...)
	}
	defer func() {
		// Add a note at the bottom about the old manifest indicator in the status.
		for _, err := range endWarnings {
			log.Warn(err)
		}
	}()

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
	isOld := ""
	if _, ok := outdatedVMs[vm.Name]; ok {
		isOld = oldManifestIndicator
	}

	if vm.Running() {
		return fmt.Sprintf("%sUp %s", isOld, vm.Status.StartTime)
	}

	return isOld + "Stopped"
}

// fetchLatestStatus fetches the current status the VMs, updates the VM status
// in memory and returns a list of outdated VMs.
func fetchLatestStatus(vms []*api.VM) (outdatedVMs map[string]bool, errList []error) {
	outdatedVMs = map[string]bool{}
	errList = []error{}

	// Container runtime clients. These clients are lazy initialized based on
	// the VM's runtime.
	var containerdClient, dockerClient runtime.Interface

	// Iterate through the VMs, fetching the actual status from the runtime.
	for _, vm := range vms {
		// Skip VMs with no runtime info or no runtime ID.
		if vm.Status.Runtime == nil || vm.Status.Runtime.ID == "" {
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
				var err error
				containerdClient, err = containerdruntime.GetContainerdClient()
				if err != nil {
					errList = append(errList, err)
					return
				}
			}
			vmRuntime = containerdClient
		case runtime.RuntimeDocker:
			if dockerClient == nil {
				var err error
				dockerClient, err = dockerruntime.GetDockerClient()
				if err != nil {
					errList = append(errList, err)
					return
				}
			}
			vmRuntime = dockerClient
		default:
			// Skip VMs with unknown runtime
			continue
		}

		// Inspect the VM container using the runtime client.
		ir, inspectErr := vmRuntime.InspectContainer(containerID)
		if inspectErr != nil {
			// Failed to get the container status. Latest status can't be
			// confirmed.
			outdatedVMs[vm.Name] = true
			errList = append(errList, errors.Wrapf(inspectErr, "failed to inspect container for VM %s", containerID))
			continue
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
