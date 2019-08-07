package main

import (
	"os"

	"github.com/weaveworks/ignite/cmd/ignite/cmd"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/providers/ignite"
	"github.com/weaveworks/ignite/pkg/util"
)

func main() {
	if err := Run(); err != nil {
		os.Exit(1)
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

	// Populate the providers
	util.GenericCheckErr(providers.Populate(ignite.Providers))

	c := cmd.NewIgniteCommand(os.Stdin, os.Stdout, os.Stderr)
	return c.Execute()
}
