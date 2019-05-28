package cmd

import "github.com/spf13/pflag"

// This file contains a collection of common flag adders
func addNameFlag(fs *pflag.FlagSet, name *string) {
	fs.StringVarP(name, "name", "n", *name, "Specify the name")
}

func addInteractiveFlag(fs *pflag.FlagSet, interactive *bool) {
	fs.BoolVarP(interactive, "interactive", "i", *interactive, "Attach to the VM after starting")
}
