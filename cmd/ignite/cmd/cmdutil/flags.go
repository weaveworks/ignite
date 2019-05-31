package cmdutil

import "github.com/spf13/pflag"

// This file contains a collection of common flag adders
func AddNameFlag(fs *pflag.FlagSet, name *string) {
	fs.StringVarP(name, "name", "n", *name, "Specify the name")
}

func AddInteractiveFlag(fs *pflag.FlagSet, interactive *bool) {
	fs.BoolVarP(interactive, "interactive", "i", *interactive, "Attach to the VM after starting")
}

func AddImportKernelFlags(fs *pflag.FlagSet, importKernel *bool, kernelName *string) {
	fs.BoolVarP(importKernel, "kernel", "k", *importKernel, "Import a new kernel from the image (if exists)")
	fs.StringVarP(kernelName, "kernel-name", "r", *kernelName, "Name the newly imported kernel (if -k is specified)")
}
