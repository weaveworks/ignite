package vmmd

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type PortMappings map[uint64]uint64

var _ fmt.Stringer = &PortMappings{}

func (pm *PortMappings) String() string {
	var sb strings.Builder
	var index int

	for hostPort, vmPort := range *pm {
		sb.WriteString(fmt.Sprintf("0.0.0.0:%d->%d", hostPort, vmPort))

		index++
		if index < len(*pm) {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}

func (md *VMMetadata) NewPortMappings(input []string) error {
	result := map[uint64]uint64{}

	for _, portMapping := range input {
		ports := strings.Split(portMapping, ":")
		if len(ports) != 2 {
			return fmt.Errorf("port mappings must be of form <host port>:<VM port>")
		}

		hostPort, err := strconv.ParseUint(ports[0], 10, 64)
		if err != nil {
			return err
		}

		vmPort, err := strconv.ParseUint(ports[1], 10, 64)
		if err != nil {
			return err
		}

		if _, ok := result[hostPort]; ok {
			return fmt.Errorf("cannot use a port on the host twice")
		}

		result[hostPort] = vmPort
	}

	md.VMOD().PortMappings = result
	return nil
}

type IPAddrs []net.IP

var _ fmt.Stringer = &IPAddrs{}

func (ip *IPAddrs) String() string {
	var sb strings.Builder

	for i, ipAddr := range *ip {
		sb.WriteString(fmt.Sprintf("%s", ipAddr))

		if i+1 < len(*ip) {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}
