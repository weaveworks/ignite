package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra/doc"
	ignitecmd "github.com/weaveworks/ignite/cmd/ignite/cmd"
	ignitedcmd "github.com/weaveworks/ignite/cmd/ignited/cmd"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/providers/ignite"
)

func main() {
	if err := providers.Populate(ignite.Providers); err != nil {
		log.Fatal(err)
	}
	ignite := ignitecmd.NewIgniteCommand(os.Stdin, os.Stdout, os.Stderr)
	if err := doc.GenMarkdownTree(ignite, "./docs/cli"); err != nil {
		log.Fatal(err)
	}
	ignited := ignitedcmd.NewIgnitedCommand(os.Stdin, os.Stdout, os.Stderr)
	if err := doc.GenMarkdownTree(ignited, "./docs/cli"); err != nil {
		log.Fatal(err)
	}
	if output, err := exec.Command("/bin/bash", "-c", `sed -e "/Auto generated/d" -i docs/cli/*.md`).CombinedOutput(); err != nil {
		log.Fatal(string(output), err)
	}
}
