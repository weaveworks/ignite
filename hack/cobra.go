package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra/doc"
	"github.com/weaveworks/ignite/cmd/ignite/cmd"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/providers/ignite"
)

func main() {
	if err := providers.Populate(ignite.Providers); err != nil {
		log.Fatal(err)
	}
	ignite := cmd.NewIgniteCommand(os.Stdin, os.Stdout, os.Stderr)
	if err := doc.GenMarkdownTree(ignite, "./docs/cli/ignite"); err != nil {
		log.Fatal(err)
	}
	if output, err := exec.Command("/bin/bash", "-c", `sed -e "/Auto generated/d" -i docs/cli/ignite/*.md`).CombinedOutput(); err != nil {
		log.Fatal(string(output), err)
	}
}
