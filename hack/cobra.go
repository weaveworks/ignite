package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra/doc"
	"github.com/weaveworks/ignite/cmd/ignite/cmd"
)

func main() {
	ignite := cmd.NewIgniteCommand(os.Stdin, os.Stdout, os.Stderr)
	if err := doc.GenMarkdownTree(ignite, "./docs/cli"); err != nil {
		log.Fatal(err)
	}
	if output, err := exec.Command("/bin/bash", "-c", `sed -e "/Auto generated/d" -i docs/cli/*.md`).CombinedOutput(); err != nil {
		log.Fatal(string(output), err)
	}
	
}
