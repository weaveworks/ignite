package main

import (
	"fmt"
	"net"
	"os"
	"path"

	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/container"
	"github.com/weaveworks/ignite/pkg/container/prometheus"
	"github.com/weaveworks/ignite/pkg/dmlegacy"
	"github.com/weaveworks/ignite/pkg/providers"
	patchutil "github.com/weaveworks/ignite/pkg/util/patch"
)

func main() {
	// Populate the providers
	cmdutil.CheckErr(providers.Populate(providers.Providers))
	RunIgniteSpawn()
}

func StartVM(co *options) error {
	// Setup networking inside of the container, return the available interfaces
	dhcpIfaces, err := container.SetupContainerNetworking()
	if err != nil {
		return fmt.Errorf("network setup failed: %v", err)
	}

	// Serve DHCP requests for those interfaces
	// This function returns the available IP addresses that are being
	// served over DHCP now
	ipAddrs, err := container.StartDHCPServers(co.vm, dhcpIfaces)
	if err != nil {
		return err
	}

	// Serve metrics over an unix socket in the VM's own directory
	metricsSocket := path.Join(co.vm.ObjectPath(), constants.PROMETHEUS_SOCKET)
	go prometheus.ServeMetrics(metricsSocket)

	// Update the VM status and IP address information
	if err := patchRunning(co.vm, ipAddrs); err != nil {
		return fmt.Errorf("failed to patch VM state: %v", err)
	}

	// Patches the VM object to set state to stopped, and clear IP addresses
	defer patchStopped(co.vm)

	// Remove the snapshot overlay post-run, which also removes the detached backing loop devices
	defer dmlegacy.DeactivateSnapshot(co.vm)

	// Remove the Prometheus socket post-run
	defer os.Remove(metricsSocket)

	// Execute Firecracker
	if err := container.ExecuteFirecracker(co.vm, dhcpIfaces); err != nil {
		return fmt.Errorf("runtime error for VM %q: %v", co.vm.GetUID(), err)
	}

	return nil
}

func patchRunning(vm *api.VM, ipAddrs []net.IP) error {
	patch, err := patchutil.Create(vm, func(obj meta.Object) error {
		patchVM := obj.(*api.VM)
		patchVM.Status.State = api.VMStateRunning
		patchVM.Status.IPAddresses = ipAddrs
		return nil
	})
	if err != nil {
		return err
	}
	// Perform the patch
	return providers.Client.VMs().Patch(vm.GetUID(), patch)
}

func patchStopped(vm *api.VM) error {
	patch, err := patchutil.Create(vm, func(obj meta.Object) error {
		patchVM := obj.(*api.VM)
		patchVM.Status.State = api.VMStateStopped
		patchVM.Status.IPAddresses = nil
		return nil
	})
	if err != nil {
		return err
	}
	// Perform the patch
	return providers.Client.VMs().Patch(vm.GetUID(), patch)
}
