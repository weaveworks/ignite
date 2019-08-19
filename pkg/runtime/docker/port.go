package docker

import (
	"fmt"
	"strconv"

	"github.com/docker/go-connections/nat"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

// portBindingsToDocker takes in portMappings and returns a nat.PortMap of the
// port bindings and a nat.PortSet of the exposed ports for the Docker client
func portBindingsToDocker(portMappings meta.PortMappings) (nat.PortMap, nat.PortSet) {
	bindings, exposed := make(nat.PortMap, len(portMappings)), make(nat.PortSet, len(portMappings))

	for _, portMapping := range portMappings {
		var hostIP string
		if portMapping.BindAddress != nil {
			hostIP = portMapping.BindAddress.String()
		}

		protocol := portMapping.Protocol
		if len(protocol) == 0 {
			// Docker uses TCP by default
			protocol = meta.ProtocolTCP
		}

		port := nat.Port(fmt.Sprintf("%d/%s", portMapping.VMPort, protocol.String()))
		exposed[port] = struct{}{}
		bindings[port] = []nat.PortBinding{
			{
				HostIP:   hostIP,
				HostPort: strconv.FormatUint(portMapping.HostPort, 10),
			},
		}
	}

	return bindings, exposed
}
