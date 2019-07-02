package cmdutil

import (
	"github.com/spf13/pflag"
	ignitemeta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

// This file contains a collection of common flag adders
func AddNameFlag(fs *pflag.FlagSet, name *string) {
	fs.StringVarP(name, "name", "n", *name, "Specify the name")
}

func AddConfigFlag(fs *pflag.FlagSet, configFile *string) {
	fs.StringVar(configFile, "config", *configFile, "Specify a path to a file with the API resources you want to pass")
}

func AddInteractiveFlag(fs *pflag.FlagSet, interactive *bool) {
	fs.BoolVarP(interactive, "interactive", "i", *interactive, "Attach to the VM after starting")
}

func AddForceFlag(fs *pflag.FlagSet, force *bool) {
	fs.BoolVarP(force, "force", "f", *force, "Force this operation. Warning, use of this mode may have unintended consequences.")
}

type SizeFlag struct {
	value ignitemeta.Size
}

func (sf *SizeFlag) Set(val string) error {
	var err error
	sf.value, err = ignitemeta.NewSizeFromString(val)
	return err
}

func (sf *SizeFlag) String() string {
	return sf.value.String()
}

func (sf *SizeFlag) Type() string {
	return "size"
}

var _ pflag.Value = &SizeFlag{}

func SizeVar(fs *pflag.FlagSet, ptr *ignitemeta.Size, name string, defVal ignitemeta.Size, usage string) {
	SizeVarP(fs, ptr, name, "", defVal, usage)
}

func SizeVarP(fs *pflag.FlagSet, ptr *ignitemeta.Size, name, shorthand string, defVal ignitemeta.Size, usage string) {
	fs.VarP(&SizeFlag{value: defVal}, name, shorthand, usage)
}
