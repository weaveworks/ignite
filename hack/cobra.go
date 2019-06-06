package main

import (
	"log"
	"os"

	"github.com/spf13/cobra/doc"
	"github.com/weaveworks/ignite/cmd/ignite/cmd"
)

func main() {
	ignite := cmd.NewIgniteCommand(os.Stdin, os.Stdout, os.Stderr)
	if err := doc.GenMarkdownTree(ignite, "./docs/cli"); err != nil {
		log.Fatal(err)
	}
}
