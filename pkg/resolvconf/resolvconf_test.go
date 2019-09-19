package resolvconf

import (
	"strings"
	"testing"

	"github.com/miekg/dns"
)

func TestBuildResolvConf(t *testing.T) {
	data := buildResolvConf(&dns.ClientConfig{
		Servers: []string{"1.1.1.1", "8.8.8.8", "9.9.9.9"},
		Search:  []string{"fire.local", "test.fire.local"},
	})
	expected := `# The following config was built by ignite:
nameserver 1.1.1.1
nameserver 8.8.8.8
nameserver 9.9.9.9
search fire.local test.fire.local
`
	if string(data) != string(expected) {
		t.Errorf(
			"\nexpected:\n%q\nactual:\n%q",
			string(expected),
			string(data),
		)
	}
}

func TestBuildResolvConfEmpty(t *testing.T) {
	data := buildResolvConf(&dns.ClientConfig{})
	expected := `# The following config was built by ignite:
`
	if string(data) != string(expected) {
		t.Errorf(
			"\nexpected:\n%q\nactual:\n%q",
			string(expected),
			string(data),
		)
	}
}

func TestIsSystemdResolved(t *testing.T) {
	type tc struct {
		cfg      dns.ClientConfig
		expected bool
	}
	testCases := []tc{
		{
			cfg: dns.ClientConfig{
				Servers: []string{"127.0.0.53"},
			},
			expected: true,
		},
		{
			cfg: dns.ClientConfig{
				Servers: []string{"::ffff:127.0.0.53"},
			},
			expected: true,
		},
		{
			cfg: dns.ClientConfig{
				Servers: []string{"::fFFf:127.0.0.53"},
			},
			expected: true,
		},
		{
			cfg: dns.ClientConfig{
				Servers: []string{"127.0.0.1", "::ffFF:127.0.0.53", "::ffff:127.0.0.53"},
			},
			expected: false,
		},
		{
			cfg: dns.ClientConfig{
				Servers: []string{"127.0.0.1", "::ffFF:127.0.0.53", "::ffff:127.0.0.53"},
			},
			expected: false,
		},
		{
			cfg: dns.ClientConfig{
				Servers: []string{"1.1.1.1"},
			},
			expected: false,
		},
		{
			cfg: dns.ClientConfig{
				Servers: []string{"127.0.0.1", "1.1.1.1"},
			},
			expected: false,
		},
		{
			cfg: dns.ClientConfig{
				Servers: []string{"127.0.0.53", "1.1.1.1", "8.8.8.8", "9.9.9.9"},
			},
			expected: false,
		},
		{
			cfg: dns.ClientConfig{
				Servers: []string{"127.0.0.1"},
			},
			expected: false,
		},
		{
			cfg: dns.ClientConfig{
				Servers: []string{"::1"},
			},
			expected: false,
		},
		{
			cfg: dns.ClientConfig{
				Servers: []string{"::53"},
			},
			expected: false,
		},
		{
			cfg: dns.ClientConfig{
				Servers: []string{},
			},
			expected: false,
		},
	}
	for i, tc := range testCases {
		res := isSystemdResolved(&tc.cfg)
		if res != tc.expected {
			t.Errorf("Case #%d expected %t, got %t: %q", i, tc.expected, res, tc.cfg.Servers)
		}
	}
}

func TestFilterLoopBack(t *testing.T) {
	type tc struct {
		cfg      dns.ClientConfig
		expected []string
	}
	testCases := []tc{
		{
			cfg: dns.ClientConfig{
				Servers: []string{},
			},
			expected: []string{},
		},
		{
			cfg: dns.ClientConfig{
				Servers: []string{"127.0.0.1", "1.1.1.1", "::ffFF:127.0.0.53", "::ffff:127.0.0.1"},
			},
			expected: []string{"1.1.1.1"},
		},
		{
			cfg: dns.ClientConfig{
				Servers: []string{"127.0.0.1", "::ffFF:127.0.0.53", "::ffff:127.0.0.53"},
			},
			expected: []string{},
		},
		{
			cfg: dns.ClientConfig{
				Servers: []string{"1.1.1.1"},
			},
			expected: []string{"1.1.1.1"},
		},
		{
			cfg: dns.ClientConfig{
				Servers: []string{"127.0.0.1", "1.1.1.1"},
			},
			expected: []string{"1.1.1.1"},
		},
		{
			cfg: dns.ClientConfig{
				Servers: []string{"127.0.0.53", "1.1.1.1", "8.8.8.8", "9.9.9.9"},
			},
			expected: []string{"1.1.1.1", "8.8.8.8", "9.9.9.9"},
		},
		{
			cfg: dns.ClientConfig{
				Servers: []string{"127.0.0.1"},
			},
			expected: []string{},
		},
		{
			cfg: dns.ClientConfig{
				Servers: []string{"::1"},
			},
			expected: []string{""},
		},
		{
			cfg: dns.ClientConfig{
				Servers: []string{"::53"},
			},
			expected: []string{"::53"},
		},
	}
	for _, tc := range testCases {
		filterLoopback(&tc.cfg)
		if strings.Join(tc.cfg.Servers, ", ") != strings.Join(tc.expected, ", ") {
			t.Errorf(
				"\nexpected:\n%v\nactual:\n%v",
				tc.expected,
				tc.cfg.Servers,
			)
		}
	}
}
