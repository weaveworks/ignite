package run

import (
	"fmt"
	"github.com/weaveworks/ignite/cmd/ignite/run/runutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/util"
	"github.com/weaveworks/ignite/pkg/version"
)

type StartFlags struct {
	PortMappings []string
	Interactive  bool
	Debug        bool
}

type startOptions struct {
	*StartFlags
	*attachOptions
}

func (sf *StartFlags) NewStartOptions(l *runutil.ResourceLoader, vmMatch string) (*startOptions, error) {
	ao, err := NewAttachOptions(l, vmMatch)
	if err != nil {
		return nil, err
	}

	// Disable running check as it takes a while for the in-container Ignite to update the state
	ao.checkRunning = false

	return &startOptions{sf, ao}, nil
}

func Start(so *startOptions) error {
	// Check if the given VM is already running
	if so.vm.Running() {
		return fmt.Errorf("VM %q is already running", so.vm.ID)
	}

	// Setup the snapshot overlay filesystem
	if err := so.vm.SetupSnapshot(); err != nil {
		return err
	}

	// Resolve the Ignite binary to be mounted inside the container
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		return err
	}
	igniteBinary, _ := filepath.Abs(path)

	dockerArgs := []string{
		"-itd",
		fmt.Sprintf("--label=ignite.name=%s", so.vm.Name.String()),
		fmt.Sprintf("--name=%s", constants.IGNITE_PREFIX+so.vm.ID),
		fmt.Sprintf("-v=%s:/ignite/ignite", igniteBinary),
		fmt.Sprintf("-v=%s:%s", constants.DATA_DIR, constants.DATA_DIR),
		fmt.Sprintf("--stop-timeout=%d", constants.STOP_TIMEOUT+constants.IGNITE_TIMEOUT),
		"--cap-add=SYS_ADMIN",          // Needed to run "dmsetup remove" inside the container
		"--cap-add=NET_ADMIN",          // Needed for removing the IP from the container's interface
		"--device=/dev/mapper/control", // This enables containerized Ignite to remove its own dm snapshot
		"--device=/dev/net/tun",        // Needed for creating TAP adapters
		"--device=/dev/kvm",            // Pass though virtualization support
		fmt.Sprintf("--device=%s", so.vm.SnapshotDev()),
	}

	dockerCmd := append(make([]string, 0, len(dockerArgs)+2), "run")

	// If we're not debugging, remove the container post-run
	if !so.Debug {
		dockerCmd = append(dockerCmd, "--rm")
	}

	ports, err := parsePortMappings(so.PortMappings)
	if err != nil {
		return err
	}

	for hostPort, vmPort := range ports {
		dockerArgs = append(dockerArgs, fmt.Sprintf("-p=%d:%d", hostPort, vmPort))
	}

	dockerArgs = append(dockerArgs, fmt.Sprintf("weaveworks/ignite:%s", version.GetFirecracker()))
	dockerArgs = append(dockerArgs, so.vm.ID)

	// Start the VM in docker
	if _, err := util.ExecuteCommand("docker", append(dockerCmd, dockerArgs...)...); err != nil {
		return fmt.Errorf("failed to start container for VM %q: %v", so.vm.ID, err)
	}

	// If starting interactively, attach after starting
	if so.Interactive {
		if err := Attach(so.attachOptions); err != nil {
			return err
		}
	} else {
		// Print the ID of the started VM
		fmt.Println(so.vm.ID)
	}

	return nil
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
