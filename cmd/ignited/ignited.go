package main

import (
	"os"
	"os/signal"
	"path"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/cmd/ignited/cmd"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/providers/ignited"
	"github.com/weaveworks/ignite/pkg/util"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-sigChan
		log.Debugf("Signal %q caught", s)
		cleanup()

		log.Debug("Program terminated normally")
		os.Exit(0)
	}()

	if err := Run(); err != nil {
		log.Debugf("Termination with error: %v", err)
		os.Exit(1)
	}
}

func cleanup() {

	var daemonSocket = path.Join(constants.DATA_DIR, constants.DAEMON_SOCKET)

	err := os.Remove(daemonSocket)
	if err == nil {
		log.Debugf("Socket %q removed successfully", daemonSocket)
	}
}

// Run runs the main cobra command of this application
func Run() error {
	// Ignite needs to run as root for now, see
	// https://github.com/weaveworks/ignite/issues/46
	// TODO: Remove this when ready
	util.GenericCheckErr(util.TestRoot())

	// Create the directories needed for running
	util.GenericCheckErr(util.CreateDirectories())

	// Preload necessary providers
	util.GenericCheckErr(providers.Populate(ignited.Preload))

	c := cmd.NewIgnitedCommand(os.Stdin, os.Stdout, os.Stderr)
	return c.Execute()
}
