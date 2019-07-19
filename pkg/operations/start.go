package operations

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"
	"time"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/dmlegacy"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/network/cni"
	"github.com/weaveworks/ignite/pkg/runtime/docker"
	"github.com/weaveworks/ignite/pkg/util"
	"github.com/weaveworks/ignite/pkg/version"
)

func StartVM(vm *vmmd.VM, debug bool) error {
	// Make sure the VM container does not exist. Don't care about the error.
	RemoveVMContainer(vm.VM)

	// Setup the snapshot overlay filesystem
	if err := dmlegacy.ActivateSnapshot(vm); err != nil {
		return err
	}

	vmDir := filepath.Join(constants.VM_DIR, vm.GetUID().String())
	kernelDir := filepath.Join(constants.KERNEL_DIR, vm.GetKernelUID().String())

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

	if vm.Spec.Network.Mode == api.NetworkModeCNI {
		dockerArgs = append(dockerArgs, "--net=none")
	}

	dockerCmd := append(make([]string, 0, len(dockerArgs)+2), "run")

	// If we're not debugging, remove the container post-run
	if !debug {
		dockerCmd = append(dockerCmd, "--rm")
	}

	// Add the port mappings to Docker
	for _, portMapping := range vm.Spec.Network.Ports {
		dockerArgs = append(dockerArgs, fmt.Sprintf("-p=%d:%d", portMapping.HostPort, portMapping.VMPort))
	}

	igniteImage := fmt.Sprintf("weaveworks/ignite:%s", version.GetIgnite().ImageTag())
	dockerArgs = append(dockerArgs, igniteImage, vm.GetUID().String())

	// Verify that the image containing ignite-spawn is pulled
	// TODO: Integrate automatic pulling into pkg/runtime
	if err := verifyPulled(igniteImage); err != nil {
		return err
	}

	// Create the VM container in docker
	// TODO: Replace all calls to the docker binary with pkg/runtime
	output, err := util.ExecuteCommand("docker", append(dockerCmd, dockerArgs...)...)
	if err != nil {
		return fmt.Errorf("failed to start container for VM %q: %v", vm.GetUID(), err)
	}

	containerID := strings.TrimSpace(output)

	if vm.Spec.Network.Mode == api.NetworkModeCNI {
		if err := setupCNINetworking(containerID); err != nil {
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

func verifyPulled(image string) error {
	client, err := docker.GetDockerClient()
	if err != nil {
		return err
	}

	if _, err = client.InspectImage(image); err != nil {
		log.Printf("Pulling image %q...", image)
		rc, err := client.PullImage(image)
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
		if _, err = client.InspectImage(image); err != nil {
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
