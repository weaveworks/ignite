package auth

import (
	"fmt"

	"github.com/containerd/containerd/remotes/docker"
	dockercliconfig "github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/credentials"
	dockercliconfigtypes "github.com/docker/cli/cli/config/types"
	log "github.com/sirupsen/logrus"
)

// NOTE: This file is based on nerdctl's dockerconfigresolver.
// Refer: https://github.com/containerd/nerdctl/blob/v0.8.1/pkg/imgutil/dockerconfigresolver/dockerconfigresolver.go

// AuthCreds is for docker.WithAuthCreds used in containerd remote resolver.
type AuthCreds func(string) (string, string, error)

// NewAuthCreds returns an AuthCreds which loads the credentials from the
// docker client config.
func NewAuthCreds(refHostname string, configPath string) (AuthCreds, string, error) {
	log.Debugf("runtime.auth: registry config dir path: %q", configPath)

	// Load does not raise an error on ENOENT
	dockerConfigFile, err := dockercliconfig.Load(configPath)
	if err != nil {
		return nil, "", err
	}

	// DefaultHost converts "docker.io" to "registry-1.docker.io",
	// which is wanted  by credFunc .
	credFuncExpectedHostname, err := docker.DefaultHost(refHostname)
	if err != nil {
		return nil, "", err
	}

	var credFunc AuthCreds
	var serverAddress string

	authConfigHostnames := []string{refHostname}
	if refHostname == "docker.io" || refHostname == "registry-1.docker.io" {
		// "docker.io" appears as ""https://index.docker.io/v1/" in ~/.docker/config.json .
		// GetAuthConfig takes the hostname part as the argument: "index.docker.io"
		authConfigHostnames = append([]string{"index.docker.io"}, refHostname)
	}

	for _, authConfigHostname := range authConfigHostnames {
		// GetAuthConfig does not raise an error on ENOENT
		ac, err := dockerConfigFile.GetAuthConfig(authConfigHostname)
		if err != nil {
			log.Errorf("cannot get auth config for authConfigHostname=%q (refHostname=%q): %v",
				authConfigHostname, refHostname, err)
		} else {
			// When refHostname is "docker.io":
			// - credFuncExpectedHostname: "registry-1.docker.io"
			// - credFuncArg:              "registry-1.docker.io"
			// - authConfigHostname:       "index.docker.io"
			// - ac.ServerAddress:         "https://index.docker.io/v1/".
			if !isAuthConfigEmpty(ac) {
				if ac.ServerAddress == "" {
					log.Warnf("failed to get ac.ServerAddress for authConfigHostname=%q (refHostname=%q)",
						authConfigHostname, refHostname)
				} else {
					acsaHostname := credentials.ConvertToHostname(ac.ServerAddress)
					if acsaHostname != authConfigHostname {
						return nil, "", fmt.Errorf("expected the hostname part of ac.ServerAddress (%q) to be authConfigHostname=%q, got %q",
							ac.ServerAddress, authConfigHostname, acsaHostname)
					}
				}

				// if ac.RegistryToken != "" {
				//     // Even containerd/CRI does not support RegistryToken as of v1.4.3,
				//     // so, nobody is actually using RegistryToken?
				//     log.Warnf("ac.RegistryToken (for %q) is not supported yet (FIXME)", authConfigHostname)
				// }

				credFunc = func(credFuncArg string) (string, string, error) {
					if credFuncArg != credFuncExpectedHostname {
						return "", "", fmt.Errorf("expected credFuncExpectedHostname=%q (refHostname=%q), got credFuncArg=%q",
							credFuncExpectedHostname, refHostname, credFuncArg)
					}
					if ac.IdentityToken != "" {
						return "", ac.IdentityToken, nil
					}
					return ac.Username, ac.Password, nil
				}
				serverAddress = ac.ServerAddress
				break
			}
		}
	}
	// credFunc can be nil here.
	return credFunc, serverAddress, nil
}

func isAuthConfigEmpty(ac dockercliconfigtypes.AuthConfig) bool {
	if ac.IdentityToken != "" || ac.Username != "" || ac.Password != "" || ac.RegistryToken != "" {
		return false
	}
	return true
}
