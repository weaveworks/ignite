package operations

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"time"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/runtime"
	"github.com/weaveworks/ignite/pkg/util"
	"github.com/weaveworks/ignite/pkg/version"
)

func StartVM(vm *vmmd.VM, debug bool) error {
	// Make sure the VM container does not exist. Don't care about the error.
	_ = RemoveVMContainer(vm.VM)

	// Setup the snapshot overlay filesystem
	if err := vm.SetupSnapshot(); err != nil {
		return err
	}

	vmDir := filepath.Join(constants.VM_DIR, vm.GetUID().String())
	kernelDir := filepath.Join(constants.KERNEL_DIR, vm.GetKernelUID().String())
	igniteImage := fmt.Sprintf("weaveworks/ignite:%s", version.GetIgnite().ImageTag())

	// Verify that the image containing ignite-spawn is pulled
	// TODO: Integrate automatic pulling into pkg/runtime
	if err := verifyPulled(igniteImage); err != nil {
		return err
	}

	config := &runtime.ContainerConfig{
		Cmd:    []string{vm.GetUID().String()},
		Labels: map[string]string{"ignite.name": vm.GetName()},
		Binds: []*runtime.Bind{
			{
				HostPath:      vmDir,
				ContainerPath: vmDir,
			},
			{
				HostPath:      kernelDir,
				ContainerPath: kernelDir,
			},
		},
		CapAdds: []string{
			"SYS_ADMIN", // Needed to run "dmsetup remove" inside the container
			"NET_ADMIN", // Needed for removing the IP from the container's interface
		},
		Devices: []string{
			"/dev/mapper/control", // This enables containerized Ignite to remove its own dm snapshot
			"/dev/net/tun",        // Needed for creating TAP adapters
			"/dev/kvm",            // Pass through virtualization support
			vm.SnapshotDev(),      // The block device to boot from
		},
		StopTimeout:  constants.STOP_TIMEOUT + constants.IGNITE_TIMEOUT,
		PortBindings: vm.Spec.Network.Ports, // Add the port mappings to Docker
	}

	// If the VM is using CNI networking, disable Docker's own implementation
	if vm.Spec.Network.Mode == api.NetworkModeCNI {
		config.NetworkMode = "none"
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

	if vm.Spec.Network.Mode == api.NetworkModeCNI {
		// TODO: Right now IP addresses aren't reclaimed when the VM is removed.
		// NetworkPlugin.RemoveContainerNetwork needs to be called when removing the VM.
		if err := providers.NetworkPlugin.SetupContainerNetwork(containerID); err != nil {
			return err
		}

		log.Printf("Networking is now handled by CNI")
	}

	log.Printf("Started Firecracker VM %q in a container with ID %q", vm.GetUID(), containerID)

	// TODO: Follow-up the container here with a defer, or dedicated goroutine. We should output
	// if it started successfully or not
	// TODO: This is temporary until we have proper communication to the container
	return waitForSpawn(vm)
}

func verifyPulled(image string) error {
	if _, err := providers.Runtime.InspectImage(image); err != nil {
		log.Printf("Pulling image %q...", image)
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
func waitForSpawn(vm *vmmd.VM) error {
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
