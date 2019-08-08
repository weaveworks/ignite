package docker

import (
	"fmt"
	"strconv"

	"github.com/docker/go-connections/nat"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

func portBindingsToPortMap(portMappings meta.PortMappings) nat.PortMap {
	portMap := make(nat.PortMap)

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

		portMap[nat.Port(fmt.Sprintf("%d/%s", portMapping.VMPort, protocol.String()))] = []nat.PortBinding{
			{
				HostIP:   hostIP,
				HostPort: strconv.FormatUint(portMapping.HostPort, 10),
			},
		}
	}

	return portMap
}
