package cmdutil

import (
	"github.com/spf13/pflag"
	ignitemeta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

// This file contains a collection of common flag adders
func AddNameFlag(fs *pflag.FlagSet, name *string) {
	fs.StringVarP(name, "name", "n", *name, "Specify the name")
}

func AddInteractiveFlag(fs *pflag.FlagSet, interactive *bool) {
	fs.BoolVarP(interactive, "interactive", "i", *interactive, "Attach to the VM after starting")
}

func AddForceFlag(fs *pflag.FlagSet, force *bool) {
	fs.BoolVarP(force, "force", "f", *force, "Force this operation. Warning, use of this mode may have unintended consequences.")
}

func AddImportKernelFlags(fs *pflag.FlagSet, kernelName *string) {
	fs.StringVarP(kernelName, "import-kernel", "k", *kernelName, "Import a new kernel from /boot/vmlinux in the image with the specified name")
}

func SizeVar(fs *pflag.FlagSet, ptr interface{}, flagName string, defaultSizeBytes int64, description string) {
	SizeVarP(fs, ptr, flagName, "", defaultSizeBytes, description)
}

func SizeVarP(fs *pflag.FlagSet, ptr interface{}, flagName, shorthand string, defaultSizeBytes int64, description string) {
	size := ignitemeta.NewSizeFromBytes(uint64(defaultSizeBytes))
	switch v := ptr.(type) {
	case *string:
		fs.StringVarP(v, flagName, shorthand, size.String(), description)
	case *int64:
		fs.Int64VarP(v, flagName, shorthand, size.Int64(), description)
	default:
		panic("invalid size flag set up")
	}
}
