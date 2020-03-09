package operations

import (
	"fmt"
	"path"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	apiruntime "github.com/weaveworks/gitops-toolkit/pkg/runtime"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/dmlegacy"
	"github.com/weaveworks/ignite/pkg/logs"
	"github.com/weaveworks/ignite/pkg/operations/lookup"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/runtime"
	"github.com/weaveworks/ignite/pkg/util"
	"github.com/weaveworks/ignite/pkg/version"
)

func StartVM(vm *api.VM, debug bool) error {
	// Inspect the VM container and remove it if it exists
	inspectResult, _ := providers.Runtime.InspectContainer(util.NewPrefixer().Prefix(vm.GetUID()))
	RemoveVMContainer(inspectResult)

	// Setup the snapshot overlay filesystem
	if err := dmlegacy.ActivateSnapshot(vm); err != nil {
		return err
	}

	kernelUID, err := lookup.KernelUIDForVM(vm, providers.Client)
	if err != nil {
		return err
	}

	vmDir := filepath.Join(constants.VM_DIR, vm.GetUID().String())
	kernelDir := filepath.Join(constants.KERNEL_DIR, kernelUID.String())

	// Set sandbox image.
	var igniteImage meta.OCIImageRef
	// If component config is found and sandbox image is set, use that as the
	// ignite sandbox image.
	if providers.ComponentConfig != nil && !providers.ComponentConfig.Spec.Sandbox.OCI.IsUnset() {
		igniteImage = providers.ComponentConfig.Spec.Sandbox.OCI
	} else {
		// Were parsing already validated data, ignore the error
		igniteImage, _ = meta.NewOCIImageRef(fmt.Sprintf("weaveworks/ignite:%s", version.GetIgnite().ImageTag()))
	}

	// Verify that the image containing ignite-spawn is pulled
	// TODO: Integrate automatic pulling into pkg/runtime
	if err := verifyPulled(igniteImage); err != nil {
		return err
	}

	config := &runtime.ContainerConfig{
		Cmd:    []string{fmt.Sprintf("--log-level=%s", logs.Logger.Level.String()), vm.GetUID().String()},
		Labels: map[string]string{"ignite.name": vm.GetName()},
		Binds: []*runtime.Bind{
			{
				HostPath:      vmDir,
				ContainerPath: vmDir,
			},
			{
				// Mount the metadata.json file specifically into the container, to a well-known place for ignite-spawn to access
				HostPath:      path.Join(vmDir, constants.METADATA),
				ContainerPath: constants.IGNITE_SPAWN_VM_FILE_PATH,
			},
			{
				// Mount the vmlinux file specifically into the container, to a well-known place for ignite-spawn to access
				HostPath:      path.Join(kernelDir, constants.KERNEL_FILE),
				ContainerPath: constants.IGNITE_SPAWN_VMLINUX_FILE_PATH,
			},
		},
		CapAdds: []string{
			"SYS_ADMIN", // Needed to run "dmsetup remove" inside the container
			"NET_ADMIN", // Needed for removing the IP from the container's interface
		},
		Devices: []*runtime.Bind{
			runtime.BindBoth("/dev/mapper/control"), // This enables containerized Ignite to remove its own dm snapshot
			runtime.BindBoth("/dev/net/tun"),        // Needed for creating TAP adapters
			runtime.BindBoth("/dev/kvm"),            // Pass through virtualization support
			runtime.BindBoth(vm.SnapshotDev()),      // The block device to boot from
		},
		StopTimeout:  constants.STOP_TIMEOUT + constants.IGNITE_TIMEOUT,
		PortBindings: vm.Spec.Network.Ports, // Add the port mappings to Docker
	}

	// Add the volumes to the container devices
	for _, volume := range vm.Spec.Storage.Volumes {
		if volume.BlockDevice == nil {
			continue // Skip all non block device volumes for now
		}

		config.Devices = append(config.Devices, &runtime.Bind{
			HostPath:      volume.BlockDevice.Path,
			ContainerPath: path.Join(constants.IGNITE_SPAWN_VOLUME_DIR, volume.Name),
		})
	}

	// Prepare the networking for the container, for the given network plugin
	if err := providers.NetworkPlugin.PrepareContainerSpec(config); err != nil {
		return err
	}

	// If we're not debugging, remove the container post-run
	if !debug {
		config.AutoRemove = true
	}

	// Run the VM container in Docker
	containerID, err := providers.Runtime.RunContainer(igniteImage, config, util.NewPrefixer().Prefix(vm.GetUID()), vm.GetUID().String())
	if err != nil {
		return fmt.Errorf("failed to start container for VM %q: %v", vm.GetUID(), err)
	}

	// Set up the networking
	result, err := providers.NetworkPlugin.SetupContainerNetwork(containerID, vm.Spec.Network.Ports...)
	if err != nil {
		return err
	}

	if !logs.Quiet {
		log.Infof("Networking is handled by %q", providers.NetworkPlugin.Name())
		log.Infof("Started Firecracker VM %q in a container with ID %q", vm.GetUID(), containerID)
	}

	// TODO: Follow-up the container here with a defer, or dedicated goroutine. We should output
	// if it started successfully or not
	// TODO: This is temporary until we have proper communication to the container
	if err := waitForSpawn(vm); err != nil {
		return err
	}

	// Set the container ID for the VM
	vm.Status.Runtime = &api.Runtime{ID: containerID}

	// Set the start time for the VM
	startTime := apiruntime.Timestamp()
	vm.Status.StartTime = &startTime

	// Append the runtime IP address of the VM to its state
	for _, addr := range result.Addresses {
		vm.Status.IPAddresses = append(vm.Status.IPAddresses, addr.IP)
	}

	// Set the VM's status to running
	vm.Status.Running = true

	// Write the state changes
	return providers.Client.VMs().Set(vm)
}

// verifyPulled pulls the ignite-spawn image if it's not present
func verifyPulled(image meta.OCIImageRef) error {
	if _, err := providers.Runtime.InspectImage(image); err != nil {
		log.Infof("Pulling image %q...", image)
		if err = providers.Runtime.PullImage(image); err != nil {
			return err
		}

		// Verify the image was pulled
		if _, err = providers.Runtime.InspectImage(image); err != nil {
			return err
		}
	}

	return nil
}

// TODO: This check for the Prometheus socket file is temporary
// until we get a proper ignite <-> ignite-spawn communication channel
func waitForSpawn(vm *api.VM) error {
	const timeout = 10 * time.Second
	const checkInterval = 100 * time.Millisecond

	startTime := time.Now()
	for time.Now().Sub(startTime) < timeout {
		time.Sleep(checkInterval)

		if util.FileExists(path.Join(vm.ObjectPath(), constants.PROMETHEUS_SOCKET)) {
			return nil
		}
	}

	return fmt.Errorf("timeout waiting for ignite-spawn startup")
}
