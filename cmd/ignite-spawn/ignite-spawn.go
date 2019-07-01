package main

import (
	"fmt"
	"os"

	"github.com/weaveworks/ignite/pkg/container"
	"github.com/weaveworks/ignite/pkg/metadata/loader"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
)

func main() {
	if err := Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// Run runs the main cobra command of this application
func Run() error {
	if len(os.Args) != 1 {
		fmt.Printf("Usage: ignite-spawn [VM ID]")
		os.Exit(0)
	}

	vmID := os.Args[0]
	opts, err := NewOptions(loader.NewResLoader(), vmID)
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

	// VM state handling
	if err := co.vm.SetState(vmmd.Running); err != nil {
		return fmt.Errorf("failed to update VM state: %v", err)
	}
	defer co.vm.SetState(vmmd.Stopped) // Performs a save, all other metadata-modifying defers need to be after this

	// Remove the snapshot overlay post-run, which also removes the detached backing loop devices
	defer co.vm.RemoveSnapshot()

	// Remove the IP addresses post-run
	defer co.vm.ClearIPAddresses()

	// Remove the port mappings post-run
	defer co.vm.ClearPortMappings()

	// Execute Firecracker
	if err := container.ExecuteFirecracker(co.vm, dhcpIfaces); err != nil {
		return fmt.Errorf("runtime error for VM %q: %v", co.vm.ID, err)
	}

	return nil
}
