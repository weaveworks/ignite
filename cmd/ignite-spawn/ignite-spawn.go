package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/container"
	"github.com/weaveworks/ignite/pkg/dmlegacy"
	"github.com/weaveworks/ignite/pkg/prometheus"
	"github.com/weaveworks/ignite/pkg/util"
	patchutil "github.com/weaveworks/libgitops/pkg/util/patch"
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

	// Explicitly set the GVK on this object
	vm.SetGroupVersionKind(api.SchemeGroupVersion.WithKind(api.KindVM.Title()))
	return vm, nil
}

func StartVM(vm *api.VM) (err error) {

	envVars := parseEnvVars(os.Environ())

	// Setup networking inside of the container, return the available interfaces
	fcIfaces, dhcpIfaces, err := container.SetupContainerNetworking(envVars)
	if err != nil {
		return fmt.Errorf("network setup failed: %v", err)
	}

	// Serve DHCP requests for those interfaces
	// This function returns the available IP addresses that are being
	// served over DHCP now
	if err = container.StartDHCPServers(vm, dhcpIfaces); err != nil {
		return
	}

	// Serve metrics over an unix socket in the VM's own directory
	metricsSocket := path.Join(vm.ObjectPath(), constants.PROMETHEUS_SOCKET)
	serveMetrics(metricsSocket)

	// Patches the VM object to set state to stopped, and clear IP addresses
	defer util.DeferErr(&err, func() error { return patchStopped(vm) })

	// Remove the snapshot overlay post-run, which also removes the detached backing loop devices
	defer util.DeferErr(&err, func() error { return dmlegacy.DeactivateSnapshot(vm) })

	// Remove the Prometheus socket post-run
	defer util.DeferErr(&err, func() error { return os.Remove(metricsSocket) })

	// Execute Firecracker
	if err = container.ExecuteFirecracker(vm, fcIfaces); err != nil {
		return fmt.Errorf("runtime error for VM %q: %v", vm.GetUID(), err)
	}

	return
}

func serveMetrics(metricsSocket string) {
	go func() {
		// Create a new registry and http.Server. Don't register custom metrics to the registry quite yet.
		_, server := prometheus.New()
		if err := prometheus.ServeOnSocket(server, metricsSocket); err != nil {
			log.Errorf("prometheus server was stopped with error: %v", err)
		}
	}()
}

// TODO: Get rid of this with the daemon architecture
func patchStopped(vm *api.VM) error {
	/*
		Perform a static patch, setting the following:
		vm.status.running = false
		vm.status.ipAddresses = nil
		vm.status.runtime = nil
		vm.status.startTime = nil
	*/

	patch := []byte(`{"status":{"running":false,"network":null,"runtime":null,"startTime":null}}`)
	return patchutil.NewPatcher(scheme.Serializer).ApplyOnFile(constants.IGNITE_SPAWN_VM_FILE_PATH, patch, vm.GroupVersionKind())
}

func parseEnvVars(vars []string) map[string]string {
	result := make(map[string]string)

	for _, arg := range vars {
		parts := strings.Split(arg, "=")
		if len(parts) != 2 {
			continue
		}

		result[parts[0]] = parts[1]
	}

	return result
}
