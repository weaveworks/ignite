package main

import (
	"os"

	"github.com/weaveworks/ignite/cmd/ignite/cmd"
	"github.com/weaveworks/ignite/pkg/providers"
)

func main() {
	if err := Run(); err != nil {
		os.Exit(1)
	}
}

// Run runs the main cobra command of this application
func Run() error {
	// Populate the providers
	if err := providers.Populate(providers.Providers); err != nil {
		return err
	}

	c := cmd.NewIgniteCommand(os.Stdin, os.Stdout, os.Stderr)
	return c.Execute()
}
