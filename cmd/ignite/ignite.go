package main

import (
	"fmt"
	"os"

	"github.com/weaveworks/ignite/cmd/ignite/cmd"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/providers/ignite"
)

func main() {
	if err := Run(); err != nil {
		// TODO: This just duplicates cobra's errors
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// Run runs the main cobra command of this application
func Run() error {
	// Populate the providers
	if err := providers.Populate(ignite.Providers); err != nil {
		return err
	}

	c := cmd.NewIgniteCommand(os.Stdin, os.Stdout, os.Stderr)
	return c.Execute()
}
