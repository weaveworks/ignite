// Copyright 2020 CNI authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

// Assigned Internet Protocol Numbers
// https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml
const (
	PROTOCOL_TCP  = 6
	PROTOCOL_UDP  = 17
	PROTOCOL_SCTP = 132
)

// getNetlinkFamily returns the Netlink IP family constant
func getNetlinkFamily(isIPv6 bool) netlink.InetFamily {
	if isIPv6 {
		return unix.AF_INET6
	}
	return unix.AF_INET
}

// DeleteConntrackEntriesForDstIP delete the conntrack entries for the connections
// specified by the given destination IP and protocol
func DeleteConntrackEntriesForDstIP(dstIP string, protocol uint8) error {
	ip := net.ParseIP(dstIP)
	if ip == nil {
		return fmt.Errorf("error deleting connection tracking state, bad IP %s", ip)
	}
	family := getNetlinkFamily(ip.To4() == nil)

	filter := &netlink.ConntrackFilter{}
	filter.AddIP(netlink.ConntrackOrigDstIP, ip)
	filter.AddProtocol(protocol)

	_, err := netlink.ConntrackDeleteFilter(netlink.ConntrackTable, family, filter)
	if err != nil {
		return fmt.Errorf("error deleting connection tracking state for protocol: %d IP: %s, error: %v", protocol, ip, err)
	}
	return nil
}

// DeleteConntrackEntriesForDstPort delete the conntrack entries for the connections specified
// by the given destination port, protocol and IP family
func DeleteConntrackEntriesForDstPort(port uint16, protocol uint8, family netlink.InetFamily) error {
	filter := &netlink.ConntrackFilter{}
	filter.AddPort(netlink.ConntrackOrigDstPort, port)
	filter.AddProtocol(protocol)

	_, err := netlink.ConntrackDeleteFilter(netlink.ConntrackTable, family, filter)
	if err != nil {
		return fmt.Errorf("error deleting connection tracking state for protocol: %d Port: %d, error: %v", protocol, port, err)
	}
	return nil
}
