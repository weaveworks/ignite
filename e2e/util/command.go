package util

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"gotest.tools/assert"
)

const (
	DefaultVMImage = "weaveworks/ignite-ubuntu"
)

// Command is an ignite command execution helper. It takes a binary and the
// arguments to run with the binary. It provides chaining methods to
// facilitate easy construction of the command.
type Command struct {
	bin string
	T   *testing.T
	Cmd *exec.Cmd
}

// NewCommand takes a go test testing.T and path to ignite binary and returns
// an initialized Command.
func NewCommand(t *testing.T, binPath string) *Command {
	return &Command{
		T:   t,
		bin: binPath,
		Cmd: exec.Command(binPath),
	}
}

// New resets the command. This should be used to reuse an existing Command and
// pass different arguments by method chaining.
func (c *Command) New() *Command {
	// NOTE: Create a whole new instance of Command. Only assigning a new
	// exec.Command to c.Cmd results in "exec: Stdout already set" error. Reuse
	// of command is not allowed as per the os/exec docs.
	return &Command{
		T:   c.T,
		bin: c.bin,
		Cmd: exec.Command(c.bin),
	}
}

// With accepts arguments to be used with the command. It returns Command and
// supports method chaining.
func (c *Command) With(args ...string) *Command {
	c.Cmd.Args = append(c.Cmd.Args, args...)
	return c
}

// WithRuntime sets the runtime argument.
func (c *Command) WithRuntime(arg string) *Command {
	return c.With("--runtime=" + arg)
}

// WithNetwork sets the network argument.
func (c *Command) WithNetwork(arg string) *Command {
	return c.With("--network-plugin=" + arg)
}

// Dir sets the command execution directory.
func (c *Command) Dir(path string) *Command {
	c.Cmd.Dir = path
	return c
}

// PassThrough makes output from the command go to the same place as this process.
func (c *Command) PassThrough() *Command {
	c.Cmd.Stderr = os.Stderr
	c.Cmd.Stdout = os.Stdout
	return c
}

// Run executes the command and performs an error check. It results in fatal
// exit of the test if an error is encountered. In order to continue the test
// on encountering an error, call Command.Cmd.CombinedOutput() or the
// appropriate method to execute the command separately.
func (c *Command) Run() {
	c.T.Helper()
	out, err := c.Cmd.CombinedOutput()
	assert.Check(c.T, err, fmt.Sprintf("cmd: \n%q\n%s", c.Cmd, out))
	if err != nil {
		c.T.Fatalf("failed to run %q: %v", c.Cmd, err)
	}
}
