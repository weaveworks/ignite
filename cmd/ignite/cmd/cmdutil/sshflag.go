package cmdutil

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/util"
)

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
