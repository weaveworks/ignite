package run

import (
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
)

// ExecFlags contains the flags supported by the exec command.
type ExecFlags struct {
	Timeout      uint32
	IdentityFile string
	Tty          bool
}

type execOptions struct {
	*ExecFlags
	vm      *api.VM
	command []string
}

// NewExecOptions constructs and returns an execOptions.
func (ef *ExecFlags) NewExecOptions(vmMatch string, command ...string) (eo *execOptions, err error) {
	eo = &execOptions{
		ExecFlags: ef,
		command:   command,
	}

	eo.vm, err = getVMForMatch(vmMatch)
	return
}

// Exec executes command in a VM based on the provided execOptions.
func Exec(eo *execOptions) error {
	return runSSH(eo.vm, eo.IdentityFile, eo.command, eo.Tty, eo.Timeout)
}
