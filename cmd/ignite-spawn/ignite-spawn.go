package main

import (
	"fmt"
	"os"
	"path"

	"github.com/weaveworks/ignite/pkg/constants"

	log "github.com/sirupsen/logrus"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/container"
	"github.com/weaveworks/ignite/pkg/container/prometheus"
	"github.com/weaveworks/ignite/pkg/logs"
)

func main() {
	if err := Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// Run runs the main cobra command of this application
func Run() error {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: ignite-spawn [VM ID]")
		os.Exit(0)
	}

	vmID := os.Args[1]
	opts, err := NewOptions(vmID)
	if err != nil {
		return err
	}

	logs.InitLogs(log.InfoLevel)

	return StartVM(opts)
}

func StartVM(co *options) error {
	// Setup networking inside of the container, return the available interfaces
	dhcpIfaces, err := container.SetupContainerNetworking()
	if err != nil {
		return fmt.Errorf("network setup failed: %v", err)
	}
	// Serve DHCP requests for those interfaces
	if err := container.StartDHCPServers(co.vm, dhcpIfaces); err != nil {
		return err
	}

	// Serve metrics over an unix socket in the VM's own directory
	metricsSocket := path.Join(co.vm.ObjectPath(), constants.PROMETHEUS_SOCKET)
	go prometheus.ServeMetrics(metricsSocket)

	// VM state handling
	// TODO: Use a .Patch here instead
	if err := setState(co.vm, api.VMStateRunning); err != nil {
		return fmt.Errorf("failed to update VM state: %v", err)
	}
	defer setState(co.vm, api.VMStateStopped) // Performs a save, all other metadata-modifying defers need to be after this

	// Remove the snapshot overlay post-run, which also removes the detached backing loop devices
	defer co.vm.DeactivateSnapshot()

	// Remove the IP addresses post-run
	defer clearIPAddresses(co.vm)

	// Remove the Prometheus socket post-run
	defer os.Remove(metricsSocket)

	// Execute Firecracker
	if err := container.ExecuteFirecracker(co.vm, dhcpIfaces); err != nil {
		return fmt.Errorf("runtime error for VM %q: %v", co.vm.GetUID(), err)
	}

	return nil
}

func setState(vm *api.VM, s api.VMState) error {
	vm.Status.State = s

	return client.VMs().Set(vm)
}

func clearIPAddresses(vm *api.VM) {
	vm.Status.IPAddresses = nil
	// TODO: This currently relies on the ordering of a set of defers. Make this more robust in the future.
}
