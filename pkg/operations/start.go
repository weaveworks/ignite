package operations

import (
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
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
	// Remove the VM container if it exists
	if err := RemoveVMContainer(vm); err != nil {
		return err
	}

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
	igniteImage := fmt.Sprintf("weaveworks/ignite:%s", version.GetIgnite().ImageTag())

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

	networkPlugin := providers.NetworkPlugins[vm.Spec.Network.Mode.String()]

	// Prepare the networking for the container, for the given network plugin
	if err := networkPlugin.PrepareContainerSpec(config); err != nil {
		return err
	}

	// If we're not debugging, remove the container post-run
	if !debug {
		config.AutoRemove = true
	}

	// Run the VM container in Docker
	containerID, err := providers.Runtime.RunContainer(igniteImage, config, util.NewPrefixer().Prefix(vm.GetUID()))
	if err != nil {
		return fmt.Errorf("failed to start container for VM %q: %v", vm.GetUID(), err)
	}

	// Set up the networking
	result, err := networkPlugin.SetupContainerNetwork(containerID)
	if err != nil {
		return err
	}

	log.Infof("Networking is handled by %q", networkPlugin.Name())
	log.Infof("Started Firecracker VM %q in a container with ID %q", vm.GetUID(), containerID)

	// TODO: Follow-up the container here with a defer, or dedicated goroutine. We should output
	// if it started successfully or not
	// TODO: This is temporary until we have proper communication to the container
	if err := waitForSpawn(vm); err != nil {
		return err
	}

	// Set the container ID for the VM
	vm.Status.Runtime = &api.Runtime{ID: containerID}

	// Set the start time for the VM
	startTime := meta.Timestamp()
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

func verifyPulled(image string) error {
	if _, err := providers.Runtime.InspectImage(image); err != nil {
		log.Infof("Pulling image %q...", image)
		rc, err := providers.Runtime.PullImage(image)
		if err != nil {
			return err
		}

		// Don't output the pull command
		if _, err = io.Copy(ioutil.Discard, rc); err != nil {
			return err
		}

		if err = rc.Close(); err != nil {
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
