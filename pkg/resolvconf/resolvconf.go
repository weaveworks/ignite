package resolvconf

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/weaveworks/ignite/pkg/util"

	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

const (
	resolvDefault = "/etc/resolv.conf"
	resolvSystemd = "/run/systemd/resolve/resolv.conf"
)

var (
	// fallbackNameServers uses Google DNS
	fallbackNameServers = []string{
		"8.8.8.8",
		"8.8.4.4",
	}
)

// EnsureResolvConf ensures a valid, usable resolvConf is written to the filepath
func EnsureResolvConf(filepath string, perm os.FileMode) error {
	cfg, err := readDNSConfig()
	if err != nil {
		log.Warn(err)
	}

	filterLoopback(cfg)
	// Fallback to default DNS servers
	if len(cfg.Servers) == 0 {
		cfg.Servers = append(cfg.Servers, fallbackNameServers...)
	}

	data := buildResolvConf(cfg)
	err = util.WriteFileIfChanged(filepath, data, perm)
	if err != nil {
		return err
	}

	return nil
}

// readDNSConfig reads settings from /etc/resolv.conf -- if those settings indicate
// systemd-resolved is in use, it reads them again from /run/systemd/resolve/resolv.conf.
// If an error occurs, cfg will be defaulted to an empty dns.ClientConfig{}.
func readDNSConfig() (*dns.ClientConfig, error) {
	cfg, err := dns.ClientConfigFromFile(resolvDefault)
	if err != nil {
		return &dns.ClientConfig{}, fmt.Errorf("Using default DNS config: %v", err)
	}

	if isSystemdResolved(cfg) {
		systemdCfg, err := dns.ClientConfigFromFile(resolvSystemd)
		if err == nil {
			return systemdCfg, nil
		}
	}

	return cfg, nil
}

// isSystemdResolved returns whether the cfg likely indicate usage of systemd-resolved.
// This function should be used with the ClientConfig from /etc/resolv.conf.
func isSystemdResolved(cfg *dns.ClientConfig) bool {
	if len(cfg.Servers) == 1 {
		ip := net.ParseIP(cfg.Servers[0])
		if ip != nil && ip.IsLoopback() && ip[len(ip)-1] == 53 {
			return true // cfg.Servers is equivalent to ["127.0.0.53"] or ["::FFff:127.0.0.53"]
		}
	}
	return false
}

// filterLoopback removes loopback addresses from cfg.Servers since they're
// not usable as an upstream resolver for a vm
func filterLoopback(cfg *dns.ClientConfig) {
	servers := make([]string, 0)
	for _, s := range cfg.Servers {
		ip := net.ParseIP(s)
		if ip != nil && !ip.IsLoopback() {
			servers = append(servers, s)
		}
	}
	cfg.Servers = servers
}

// buildResolvConf returns the bytes for a resolv.conf representing clientConfig
func buildResolvConf(cfg *dns.ClientConfig) []byte {
	s := "# The following config was built by ignite:\n"
	if len(cfg.Servers) > 0 {
		s += "nameserver " +
			strings.Join(cfg.Servers, "\nnameserver ") +
			"\n"
	}
	if len(cfg.Search) > 0 {
		s += "search " + strings.Join(cfg.Search, " ") +
			"\n"
	}
	return []byte(s)
}
