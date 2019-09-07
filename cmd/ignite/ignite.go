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
	// Preload necessary providers
	util.GenericCheckErr(providers.Populate(ignite.Preload))

	c := cmd.NewIgniteCommand(os.Stdin, os.Stdout, os.Stderr)
	return c.Execute()
}
