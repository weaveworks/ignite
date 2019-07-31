package cmdutil

import (
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/logs"
	"golang.org/x/crypto/ssh"
)

// CheckErr is used by Ignite commands to check if the action failed
// and respond with a fatal error provided by the logger (calls os.Exit)
func CheckErr(err error) {
	switch e := err.(type) {
	case nil:
		return // Don't fail if there's no error
	case *ssh.ExitError: // In case of SSH errors, use the exit status of the remote command
		logs.Logger.ExitCode = e.ExitStatus()
	}

	log.Fatal(err)
}
