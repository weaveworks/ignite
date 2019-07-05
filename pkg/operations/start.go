package operations

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/network/cni"
	"github.com/weaveworks/ignite/pkg/runtime/docker"
	"github.com/weaveworks/ignite/pkg/util"
	"github.com/weaveworks/ignite/pkg/version"
)

const (
	NetworkModeCNI    = "cni"
	NetworkModeBridge = "bridge"
)

var NetworkModes = []string{
	NetworkModeCNI,
	NetworkModeBridge,
}

func StartVM(vm *vmmd.VM, networkMode string, debug bool) error {
	// Make sure the VM container does not exist. Don't care about the error
	RemoveVMContainer(vm.VM)

	// Setup the snapshot overlay filesystem
	if err := vm.SetupSnapshot(); err != nil {
		return err
	}

	vmDir := filepath.Join(constants.VM_DIR, vm.GetUID().String())
	kernelDir := filepath.Join(constants.KERNEL_DIR, vm.Spec.Kernel.UID.String())

	dockerArgs := []string{
		"-itd",
		fmt.Sprintf("--label=ignite.name=%s", vm.GetName()),
		fmt.Sprintf("--name=%s", constants.IGNITE_PREFIX+vm.GetUID()),
		fmt.Sprintf("--volume=%s:%s", vmDir, vmDir),
		fmt.Sprintf("--volume=%s:%s", kernelDir, kernelDir),
		fmt.Sprintf("--stop-timeout=%d", constants.STOP_TIMEOUT+constants.IGNITE_TIMEOUT),
		"--cap-add=SYS_ADMIN",          // Needed to run "dmsetup remove" inside the container
		"--cap-add=NET_ADMIN",          // Needed for removing the IP from the container's interface
		"--device=/dev/mapper/control", // This enables containerized Ignite to remove its own dm snapshot
		"--device=/dev/net/tun",        // Needed for creating TAP adapters
		"--device=/dev/kvm",            // Pass though virtualization support
		fmt.Sprintf("--device=%s", vm.SnapshotDev()),
	}

	if networkMode == NetworkModeCNI {
		dockerArgs = append(dockerArgs, "--net=none")
	}

	dockerCmd := append(make([]string, 0, len(dockerArgs)+2), "run")

	// If we're not debugging, remove the container post-run
	if debug {
		dockerCmd = append(dockerCmd, "--rm")
	}

	// Add the port mappings to Docker
	for _, portMapping := range vm.Spec.Ports {
		dockerArgs = append(dockerArgs, fmt.Sprintf("-p=%d:%d", portMapping.HostPort, portMapping.VMPort))
	}

	// Use the :dev image tag for non-release builds
	imageTag := version.GetIgnite().GitVersion
	if version.GetIgnite().GitTreeState == "dirty" {
		imageTag = "dev"
	}
	dockerArgs = append(dockerArgs, fmt.Sprintf("weaveworks/ignite:%s", imageTag))
	dockerArgs = append(dockerArgs, vm.GetUID().String())

	// Create the VM container in docker
	containerID, err := util.ExecuteCommand("docker", append(dockerCmd, dockerArgs...)...)
	if err != nil {
		return fmt.Errorf("failed to start container for VM %q: %v", vm.GetUID(), err)
	}

	if networkMode == NetworkModeCNI {
		if err := setupCNINetworking(containerID); err != nil {
			return err
		}
		log.Printf("Networking is now handled by CNI")
	}
	log.Printf("Started Firecracker VM %q in a container with ID %q", vm.GetUID(), containerID)
	// TODO: Follow-up the container here with a defer, or dedicated goroutine. We should output
	// if it started successfully or not
	return nil
}

func setupCNINetworking(containerID string) error {
	// TODO: Both the client and networkPlugin variables should be constructed once,
	// and accessible throughout the program.
	// TODO: Right now IP addresses aren't reclaimed when the VM is removed.
	// networkPlugin.RemoveContainerNetwork need to be called when removing the VM.
	client, err := docker.GetDockerClient()
	if err != nil {
		return err
	}
	networkPlugin, err := cni.GetCNINetworkPlugin(client)
	if err != nil {
		return err
	}
	return networkPlugin.SetupContainerNetwork(containerID)
}
