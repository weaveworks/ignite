package run

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"time"

	"github.com/alessio/shellescape"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/util"
	"golang.org/x/crypto/ssh"
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
	// Check if the VM is running
	if !eo.vm.Running() {
		return fmt.Errorf("VM %q is not running", eo.vm.GetUID())
	}

	// Get the IP address
	ipAddrs := eo.vm.Status.IPAddresses
	if len(ipAddrs) == 0 {
		return fmt.Errorf("VM %q has no usable IP addresses", eo.vm.GetUID())
	}

	// If an external identity file is specified, use it instead of the internal one
	privKeyFile := eo.IdentityFile
	if len(privKeyFile) == 0 {
		privKeyFile = path.Join(eo.vm.ObjectPath(), fmt.Sprintf(constants.VM_SSH_KEY_TEMPLATE, eo.vm.GetUID()))
		if !util.FileExists(privKeyFile) {
			return fmt.Errorf("no private key found for VM %q", eo.vm.GetUID())
		}
	}

	signer, err := newSignerForKey(privKeyFile)
	if err != nil {
		return fmt.Errorf("unable to create signer for private key: %v", err)
	}

	// Create an SSH client, and connect, we will use this to exec
	config := newSSHConfig(signer, eo.Timeout)
	client, err := ssh.Dial("tcp", net.JoinHostPort(ipAddrs[0].String(), "22"), config)
	if err != nil {
		return fmt.Errorf("failed to dial: %v", err)
	}

	// Run the command, DO NOT wrap this error as the caller can check for the command exit
	// code in the ssh.ExitError type
	return runSSHCommand(client, eo.Tty, eo.command)
}

func newSignerForKey(keyPath string) (ssh.Signer, error) {
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	return ssh.ParsePrivateKey(key)
}

func newSSHConfig(publicKey ssh.Signer, timeout uint32) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(publicKey),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: use ssh.FixedPublicKey instead
		Timeout:         time.Second * time.Duration(timeout),
	}
}

func runSSHCommand(client *ssh.Client, tty bool, command []string) error {
	// create a session for the command
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	if tty {
		// get a pty
		// TODO: should these be based on the host terminal?
		// TODO: should we request something other than xterm?
		// TODO: we should probably configure the terminal modes
		modes := ssh.TerminalModes{}
		if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
			return fmt.Errorf("request for pseudo terminal failed: %v", err)
		}
	}

	// Connect input / output
	// TODO: these should come from the cobra command instead of hardcoding os.Stderr etc.
	session.Stderr = os.Stderr
	session.Stdout = os.Stdout
	session.Stdin = os.Stdin

	/*
		Do not wrap this error so the caller can check for the exit code
		If the remote server does not send an exit status, an error of type *ExitMissingError is returned.
		If the command completes unsuccessfully or is interrupted by a signal, the error is of type *ExitError.
		Other error types may be returned for I/O problems.
	*/
	return session.Run(joinShellCommand(command))
}

// joinShellCommand joins command parts into a single string safe for passing to sh -c (or SSH)
func joinShellCommand(command []string) string {
	joined := command[0]
	if len(command) == 1 {
		return joined
	}
	for _, arg := range command[1:] {
		// NOTE: we need to escape / quote to ensure that
		// each component of command... is read as a single shell word
		joined += " " + shellescape.Quote(arg)
	}
	return joined
}
