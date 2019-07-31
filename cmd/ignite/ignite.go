package main

import (
	"os"

	"github.com/weaveworks/ignite/cmd/ignite/cmd"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/providers/ignite"
)

func main() {
	if err := Run(); err != nil {
		os.Exit(1)
	}
}

// Run runs the main cobra command of this application
func Run() error {
	// Populate the providers
	cmdutil.CheckErr(providers.Populate(ignite.Providers))

	c := cmd.NewIgniteCommand(os.Stdin, os.Stdout, os.Stderr)
	return c.Execute()
}
