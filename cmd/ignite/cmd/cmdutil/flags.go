package cmdutil

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/pkg/constants"
)

// This file contains a collection of common flag adders
func AddNameFlag(fs *pflag.FlagSet, name *string) {
	fs.StringVarP(name, "name", "n", *name, "Specify the name")
}

func AddIDPrefixFlag(fs *pflag.FlagSet, idPrefix *string) {
	// Note that the flag default is printed for good UX, but it's not implemented by pflag
	// We have our own defaulting ComponentConfig logic
	fs.StringVar(idPrefix, "id-prefix", "",
		fmt.Sprintf("Prefix string for system identifiers (default %v)", constants.IGNITE_PREFIX))
}

func AddConfigFlag(fs *pflag.FlagSet, configFile *string) {
	fs.StringVar(configFile, "config", *configFile, "Specify a path to a file with the API resources you want to pass")
}

func AddInteractiveFlag(fs *pflag.FlagSet, interactive *bool) {
	fs.BoolVarP(interactive, "interactive", "i", *interactive, "Attach to the VM after starting")
}

func AddForceFlag(fs *pflag.FlagSet, force *bool) {
	fs.BoolVarP(force, "force", "f", *force, "Force this operation. Warning, use of this mode may have unintended consequences.")
}

func AddSSHFlags(fs *pflag.FlagSet, identityFile *string, timeout *uint32) {
	fs.StringVarP(identityFile, "identity", "i", "", "Override the vm's default identity file")
	fs.Uint32Var(timeout, "timeout", constants.SSH_DEFAULT_TIMEOUT_SECONDS, "Timeout waiting for connection in seconds")
}

func AddRegistryConfigDirFlag(fs *pflag.FlagSet, dir *string) {
	fs.StringVar(dir, "registry-config-dir", "", "Directory containing the registry configuration (default ~/.docker/)")
}
