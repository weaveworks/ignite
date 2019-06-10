package run

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/util"
	"github.com/weaveworks/ignite/pkg/version"
)

type StartOptions struct {
	AttachOptions
	PortMappings []string
	Interactive  bool
	Debug        bool
}

func Start(so *StartOptions) (string, error) {
	// Check if the given VM is already running
	if so.VM.Running() {
		return "", fmt.Errorf("VM %q is already running", so.VM.ID)
	}

	// Setup the snapshot overlay filesystem
	if err := so.VM.SetupSnapshot(); err != nil {
		return "", err
	}

	// Resolve the Ignite binary to be mounted inside the container
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	igniteBinary, _ := filepath.Abs(path)

	vmDir := filepath.Join(constants.VM_DIR, so.VM.ID)
	kernelDir := filepath.Join(constants.KERNEL_DIR, so.VM.KernelID())

	dockerArgs := []string{
		"-itd",
		fmt.Sprintf("--label=ignite.name=%s", so.VM.Name.String()),
		fmt.Sprintf("--name=%s", constants.IGNITE_PREFIX+so.VM.ID),
		fmt.Sprintf("--volume=%s:/ignite/ignite", igniteBinary),
		fmt.Sprintf("--volume=%s:%s", vmDir, vmDir),
		fmt.Sprintf("--volume=%s:%s", kernelDir, kernelDir),
		fmt.Sprintf("--stop-timeout=%d", constants.STOP_TIMEOUT+constants.IGNITE_TIMEOUT),
		"--cap-add=SYS_ADMIN",          // Needed to run "dmsetup remove" inside the container
		"--cap-add=NET_ADMIN",          // Needed for removing the IP from the container's interface
		"--device=/dev/mapper/control", // This enables containerized Ignite to remove its own dm snapshot
		"--device=/dev/net/tun",        // Needed for creating TAP adapters
		"--device=/dev/kvm",            // Pass though virtualization support
		fmt.Sprintf("--device=%s", so.VM.SnapshotDev()),
	}

	dockerCmd := append(make([]string, 0, len(dockerArgs)+2), "run")

	// If we're not debugging, remove the container post-run
	if !so.Debug {
		dockerCmd = append(dockerCmd, "--rm")
	}

	ports, err := parsePortMappings(so.PortMappings)
	if err != nil {
		return "", err
	}

	for hostPort, vmPort := range ports {
		dockerArgs = append(dockerArgs, fmt.Sprintf("-p=%d:%d", hostPort, vmPort))
	}

	dockerArgs = append(dockerArgs, fmt.Sprintf("weaveworks/ignite:%s", version.GetFirecracker()))
	dockerArgs = append(dockerArgs, so.VM.ID)

	// Start the VM in docker
	containerID, err := util.ExecuteCommand("docker", append(dockerCmd, dockerArgs...)...)
	if err != nil {
		return "", fmt.Errorf("failed to start container for VM %q: %v", so.VM.ID, err)
	}
	log.Printf("Started Firecracker in a Docker container with ID %q", containerID)

	// If starting interactively, attach after starting
	if so.Interactive {
		return "", Attach(&so.AttachOptions)
	}
	return so.VM.ID, nil
}

func parsePortMappings(portMappings []string) (map[uint64]uint64, error) {
	result := map[uint64]uint64{}
	for _, portMapping := range portMappings {
		ports := strings.Split(portMapping, ":")
		if len(ports) != 2 {
			return nil, fmt.Errorf("invalid --ports must be of the form hostPort:vmPort")
		}
		hostPort, err := strconv.ParseUint(ports[0], 10, 64)
		if err != nil {
			return nil, err
		}
		vmPort, err := strconv.ParseUint(ports[1], 10, 64)
		if err != nil {
			return nil, err
		}
		if _, ok := result[hostPort]; ok {
			return nil, fmt.Errorf("you can't specify two hostports twice")
		}
		result[hostPort] = vmPort
	}
	return result, nil
}
