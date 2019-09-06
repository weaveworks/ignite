package containerd

import (
	"fmt"
	"os"

	"github.com/containerd/containerd/cmd/containerd/command"
	"github.com/containerd/containerd/pkg/seed"
)

func Main() {
	seed.WithTimeAndRand()
	app := command.App()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "containerd: %s\n", err)
		os.Exit(1)
	}
}
