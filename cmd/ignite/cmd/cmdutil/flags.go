package cmdutil

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/util"
)

// This file contains a collection of common flag adders
func AddNameFlag(fs *pflag.FlagSet, name *string) {
	fs.StringVarP(name, "name", "n", *name, "Specify the name")
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

type SizeFlag struct {
	value *meta.Size
}

func (sf *SizeFlag) Set(val string) error {
	var err error
	*sf.value, err = meta.NewSizeFromString(val)
	return err
}

func (sf *SizeFlag) String() string {
	return sf.value.String()
}

func (sf *SizeFlag) Type() string {
	return "size"
}

var _ pflag.Value = &SizeFlag{}

func SizeVar(fs *pflag.FlagSet, ptr *meta.Size, name, usage string) {
	SizeVarP(fs, ptr, name, "", usage)
}

func SizeVarP(fs *pflag.FlagSet, ptr *meta.Size, name, shorthand, usage string) {
	fs.VarP(&SizeFlag{value: ptr}, name, shorthand, usage)
}

type OCIImageRefFlag struct {
	value *meta.OCIImageRef
}

func (of *OCIImageRefFlag) Set(val string) error {
	var err error
	*of.value, err = meta.NewOCIImageRef(val)
	return err
}

func (of *OCIImageRefFlag) String() string {
	if of.value == nil {
		return ""
	}
	return of.value.String()
}

func (of *OCIImageRefFlag) Type() string {
	return "oci-image"
}

var _ pflag.Value = &OCIImageRefFlag{}

func OCIImageRefVar(fs *pflag.FlagSet, ptr *meta.OCIImageRef, name, usage string) {
	OCIImageRefVarP(fs, ptr, name, "", usage)
}

func OCIImageRefVarP(fs *pflag.FlagSet, ptr *meta.OCIImageRef, name, shorthand, usage string) {
	fs.VarP(&OCIImageRefFlag{value: ptr}, name, shorthand, usage)
}

type NetworkModeFlag struct {
	value *api.NetworkMode
}

func (nf *NetworkModeFlag) Set(val string) error {
	*nf.value = api.NetworkMode(val)
	return nil
}

func (nf *NetworkModeFlag) String() string {
	if nf.value == nil {
		return ""
	}
	return nf.value.String()
}

func (nf *NetworkModeFlag) Type() string {
	return "network-mode"
}

var _ pflag.Value = &NetworkModeFlag{}

func NetworkModeVar(fs *pflag.FlagSet, ptr *api.NetworkMode) {
	fs.Var(&NetworkModeFlag{value: ptr}, "net", fmt.Sprintf("Networking mode to use. Available options are: %v", api.GetNetworkModes()))
}

// SSHFlag is the pflag.Value custom flag for `ignite create --ssh`
type SSHFlag struct {
	value *api.SSH
}

var _ pflag.Value = &SSHFlag{}

func (sf *SSHFlag) Set(x string) error {
	if x == "<path>" { // only --ssh was specified, then this default "no-value" string is set
		*sf.value = api.SSH{
			Generate: true,
		}
	} else if len(x) > 0 { // some other path was set
		importKey := x
		// Always digest the public key
		if !strings.HasSuffix(importKey, ".pub") {
			importKey = fmt.Sprintf("%s.pub", importKey)
		}
		// verify the file exists
		if !util.FileExists(importKey) {
			return fmt.Errorf("invalid SSH key: %s", importKey)
		}

		// Set the SSH PublicKey field
		*sf.value = api.SSH{
			PublicKey: importKey,
		}
	}
	return nil
}

func (sf *SSHFlag) String() string {
	if sf.value == nil {
		return ""
	}

	return sf.value.PublicKey
}

func (sf *SSHFlag) Type() string {
	return ""
}

func (sf *SSHFlag) IsBoolFlag() bool {
	return true
}

func SSHVar(fs *pflag.FlagSet, ptr *api.SSH) {
	fs.Var(&SSHFlag{value: ptr}, "ssh", "Enable SSH for the VM. If <path> is given, it will be imported as the public key. If just '--ssh' is specified, a new keypair will be generated.")

	sshFlag := fs.Lookup("ssh")
	sshFlag.NoOptDefVal = "<path>"
	sshFlag.DefValue = "is unset, which disables SSH access to the VM"
}
