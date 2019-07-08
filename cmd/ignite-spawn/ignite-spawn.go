package main

import (
	"fmt"
	"os"
	"path"

	"github.com/weaveworks/ignite/pkg/container/prometheus"

	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"

	"github.com/weaveworks/ignite/pkg/container"
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
	// TODO: when restarting a VM using `start`, we get a panic:
	// panic: listen unix /var/lib/firecracker/vm/6a2b6ebafcb0e75c/prometheus.sock: bind: address already in use
	metricsSocket := path.Join(co.vm.ObjectPath(), "prometheus.sock")
	go prometheus.ServeMetrics(metricsSocket)

	// VM state handling
	if err := co.vm.SetState(api.VMStateRunning); err != nil {
		return fmt.Errorf("failed to update VM state: %v", err)
	}
	defer co.vm.SetState(api.VMStateStopped) // Performs a save, all other metadata-modifying defers need to be after this

	// Remove the snapshot overlay post-run, which also removes the detached backing loop devices
	defer co.vm.RemoveSnapshot()

	// Remove the IP addresses post-run
	defer co.vm.ClearIPAddresses()

	// Remove the port mappings post-run
	defer co.vm.ClearPortMappings()

	// Execute Firecracker
	if err := container.ExecuteFirecracker(co.vm, dhcpIfaces); err != nil {
		return fmt.Errorf("runtime error for VM %q: %v", co.vm.GetUID(), err)
	}

	return nil
}
