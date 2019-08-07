package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
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

	cmds := map[string]*cobra.Command{
		"ignite":  ignitecmd.NewIgniteCommand(os.Stdin, os.Stdout, os.Stderr),
		"ignited": ignitedcmd.NewIgnitedCommand(os.Stdin, os.Stdout, os.Stderr),
	}

	for name, cmd := range cmds {
		if err := doc.GenMarkdownTree(cmd, fmt.Sprintf("./docs/cli/%s", name)); err != nil {
			log.Fatal(err)
		}
		sedCmd := fmt.Sprintf(`sed -e "/Auto generated/d" -i docs/cli/%s/*.md`, name)
		if output, err := exec.Command("/bin/bash", "-c", sedCmd).CombinedOutput(); err != nil {
			log.Fatal(string(output), err)
		}
	}
}
