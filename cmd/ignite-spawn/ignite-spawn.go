package main

import (
	"fmt"
	"net"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/container"
	dmcleanup "github.com/weaveworks/ignite/pkg/dmlegacy/cleanup"
	"github.com/weaveworks/ignite/pkg/prometheus"
	patchutil "github.com/weaveworks/ignite/pkg/util/patch"
)

func main() {
	// Populate the providers
	RunIgniteSpawn()
}

func decodeVM(vmID string) (*api.VM, error) {
	filePath := constants.IGNITE_SPAWN_VM_FILE_PATH
	obj, err := scheme.Serializer.DecodeFile(filePath, true)
	if err != nil {
		return nil, err
	}
	vm, ok := obj.(*api.VM)
	if !ok {
		return nil, fmt.Errorf("object couldn't be converted to VM")
	}
	// Explicitely set the GVK on this object
	vm.SetGroupVersionKind(api.SchemeGroupVersion.WithKind(api.KindVM.Title()))
	return vm, nil
}

func StartVM(vm *api.VM) error {
	// Setup networking inside of the container, return the available interfaces
	dhcpIfaces, err := container.SetupContainerNetworking()
	if err != nil {
		return fmt.Errorf("network setup failed: %v", err)
	}

	// Serve DHCP requests for those interfaces
	// This function returns the available IP addresses that are being
	// served over DHCP now
	ipAddrs, err := container.StartDHCPServers(vm, dhcpIfaces)
	if err != nil {
		return err
	}

	// Serve metrics over an unix socket in the VM's own directory
	metricsSocket := path.Join(vm.ObjectPath(), constants.PROMETHEUS_SOCKET)
	serveMetrics(metricsSocket)

	// Update the VM status and IP address information
	if err := patchRunning(vm, ipAddrs); err != nil {
		return fmt.Errorf("failed to patch VM state: %v", err)
	}

	// Patches the VM object to set state to stopped, and clear IP addresses
	defer patchStopped(vm)

	// Remove the snapshot overlay post-run, which also removes the detached backing loop devices
	defer dmcleanup.DeactivateSnapshot(vm)

	// Remove the Prometheus socket post-run
	defer os.Remove(metricsSocket)

	// Execute Firecracker
	if err := container.ExecuteFirecracker(vm, dhcpIfaces); err != nil {
		return fmt.Errorf("runtime error for VM %q: %v", vm.GetUID(), err)
	}

	return nil
}

func serveMetrics(metricsSocket string) {
	go func() {
		// create a new registry and http.Server. don't register custom metrics to the registry quite yet
		_, server := prometheus.New()
		if err := prometheus.ServeOnSocket(server, metricsSocket); err != nil {
			log.Errorf("prometheus server was stopped with error: %v", err)
		}
	}()
}

func patchRunning(vm *api.VM, ipAddrs []net.IP) error {
	return patchVM(vm, func(patchVM *api.VM) error {
		patchVM.Status.Running = true
		patchVM.Status.IPAddresses = ipAddrs
		return nil
	})
}

func patchStopped(vm *api.VM) error {
	return patchVM(vm, func(patchVM *api.VM) error {
		patchVM.Status.Running = false
		patchVM.Status.IPAddresses = nil
		return nil
	})
}

func patchVM(vm *api.VM, fn func(*api.VM) error) error {
	patch, err := patchutil.Create(vm, func(obj meta.Object) error {
		patchVM := obj.(*api.VM)
		return fn(patchVM)
	})
	if err != nil {
		return err
	}
	// Perform the patch
	return patchutil.ApplyOnFile(constants.IGNITE_SPAWN_VM_FILE_PATH, patch, vm.GroupVersionKind())
}
