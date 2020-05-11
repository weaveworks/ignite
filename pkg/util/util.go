package util

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/goombaio/namegenerator"
	"github.com/rjeczalik/notify"
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/constants"
	//"github.com/containerd/fifo"
	//"golang.org/x/net/context"
)

// GenericCheckErr is used by the commands to check if the action failed
// and respond with a fatal error provided by the logger (calls os.Exit)
// Ignite has its own, more detailed implementation of this in cmdutil
func GenericCheckErr(err error) {
	switch err.(type) {
	case nil:
		return // Don't fail if there's no error
	}

	log.Fatal(err)
}

func ExecuteCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command %q exited with %q: %v", cmd.Args, out, err)
	}

	return string(bytes.TrimSpace(out)), nil
}

func ExecForeground(command string, args ...string) (int, error) {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	cmdArgs := strings.Join(cmd.Args, " ")

	var cmdErr error
	var exitCode int

	if err != nil {
		cmdErr = fmt.Errorf("external command %q exited with an error: %v", cmdArgs, err)

		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			cmdErr = fmt.Errorf("failed to get exit code for external command %q", cmdArgs)
		}
	}

	return exitCode, cmdErr
}

func IsEmptyString(input string) bool {
	return len(strings.TrimSpace(input)) == 0
}

// Fills the given string slice with unique MAC addresses
func NewMAC(buffer *[]string) error {
	var mac string
	var macBytes []byte

	for {
		if len(*buffer) == cap(*buffer) {
			break
		}

		macBytes = make([]byte, 6)
		if _, err := rand.Read(macBytes); err != nil {
			return fmt.Errorf("failed to generate MAC: %v", err)
		}

		// Set local bit, ensure unicast address
		macBytes[0] = (macBytes[0] | 2) & 0xfe

		// Convert the byte slice to a string literally
		mac = fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", macBytes[0], macBytes[1], macBytes[2], macBytes[3], macBytes[4], macBytes[5])

		// If the generated MAC is unique break the generator loop
		unique := true
		for _, testMac := range *buffer {
			if mac == testMac {
				unique = false
				break
			}
		}

		// Generate a new MAC if it's not unique
		if !unique {
			continue
		}

		*buffer = append(*buffer, mac)
	}

	return nil
}

func NewUID() (s string, err error) {
	b := make([]byte, constants.IGNITE_UID_LENGTH/2)
	if _, err = rand.Read(b); err == nil {
		// Convert the byte slice to a string literally
		s = fmt.Sprintf("%x", b)
	}

	return
}

func RandomName() string {
	return namegenerator.NewNameGenerator(time.Now().UTC().UnixNano()).Generate()
}

func TestRoot() error {
	if syscall.Getuid() == 0 {
		return nil
	}
	return fmt.Errorf("This program needs to run as root.")
}

// This is a light weight handler to capture
// errors during detach that would otherwise
// be silently ignored. Pass a pointer to the
// error to be returned and the function to run.
// TODO: Replace all ignored defers with this
func DeferErr(err *error, f func() error) {
	if err == nil {
		panic("nil pointer given to DeferErr")
	}

	// If the given function returned an error
	// and there isn't already an error to be
	// returned, assign the function's error
	if fErr := f(); fErr != nil && *err == nil {
		*err = fErr
	}
}

type Prefixer struct {
	prefix    string
	separator string
}

func NewPrefixer() *Prefixer {
	return &Prefixer{
		prefix:    constants.IGNITE_PREFIX,
		separator: "-",
	}
}

func (p *Prefixer) Prefix(input ...interface{}) string {
	if len(input) > 0 {
		s := make([]string, 0, len(input))

		for _, data := range input {
			s = append(s, fmt.Sprintf("%v", data))
		}

		p.prefix += p.separator + strings.Join(s, p.separator)
	}

	return p.prefix
}

// Watch the fifo (named-pipe) file at pipePath and copy any
// writes to the fifo buffer to a persitent file at filePath
// Credit: https://stackoverflow.com/questions/45443414/read-continuously-from-a-named-pipe
func NamedPipeWatcher(pipePath string, filePath string) {
	var p, f *os.File
	var err error
	var e notify.EventInfo

	syscall.Mkfifo(pipePath, 0666)

	if p, err = os.Open(pipePath); os.IsNotExist(err) {
		log.Fatalf("Named pipe '%s' does not exist", pipePath)
	} else if os.IsPermission(err) {
		log.Fatalf("Insufficient permissions to read named pipe '%s': %s", pipePath, err)
	} else if err != nil {
		log.Fatalf("Error while opening named pipe '%s': %s", pipePath, err)
	}

	defer p.Close()

	// Do the same for the output file
	if f, err = os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666); os.IsNotExist(err) {
		log.Fatalf("File '%s' does not exist", filePath)
		return
	} else if os.IsPermission(err) {
		log.Fatalf("Insufficient permissions to open/create file '%s' for appending: %s", filePath, err)
		return
	} else if err != nil {
		log.Fatalf("Error while opening file '%s' for writing: %err", filePath, err)
		return
	}

	defer f.Close()

	// Here is where it happens. We create a buffered channel for events which might happen
	// on the file. The reason why we make it buffered to the number of expected concurrent writers
	// is that if all writers would (theoretically) write at once or at least pretty close
	// to each other, it might happen that we loose event. This is due to the underlying implementations
	// not because of go.
	c := make(chan notify.EventInfo, constants.MAX_CONCURRENT_WRITERS)

	// Here we tell notify to watch out named pipe for events, Write and Remove events
	// specifically. We watch for remove events, too, as this removes the file handle we
	// read from, making reads impossible
	notify.Watch(pipePath, c, notify.Write|notify.Remove)

loop:
	for {
		// ...waiting for an event to be passed.
		e = <-c

		switch e.Event() {
		case notify.Write:
			// If it a a write event, we copy the content of the named pipe to
			// our output file and wait for the next event to happen.
			// Note that this is idempotent: Even if we have huge writes by multiple
			// writers on the named pipe, the first call to Copy will copy the contents.
			// The time to copy that data may well be longer than it takes to generate the events.
			// However, subsequent calls may copy nothing, but that does not do any harm.
			io.Copy(f, p)
		case notify.Remove:
			// Some user or process has removed the named pipe,
			// so we have nothing left to read from.
			break loop
		}
	}
}
