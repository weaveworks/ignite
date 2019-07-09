package cmdutil

import (
	"github.com/spf13/pflag"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
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
	value meta.Size
}

func (sf *SizeFlag) Set(val string) error {
	var err error
	sf.value, err = meta.NewSizeFromString(val)
	return err
}

func (sf *SizeFlag) String() string {
	return sf.value.String()
}

func (sf *SizeFlag) Type() string {
	return "size"
}

var _ pflag.Value = &SizeFlag{}

func SizeVar(fs *pflag.FlagSet, ptr *meta.Size, name string, defVal meta.Size, usage string) {
	SizeVarP(fs, ptr, name, "", defVal, usage)
}

func SizeVarP(fs *pflag.FlagSet, ptr *meta.Size, name, shorthand string, defVal meta.Size, usage string) {
	fs.VarP(&SizeFlag{value: defVal}, name, shorthand, usage)
}

type OCIImageRefFlag struct {
	value meta.OCIImageRef
}

func (of *OCIImageRefFlag) Set(val string) error {
	var err error
	of.value, err = meta.NewOCIImageRef(val)
	return err
}

func (of *OCIImageRefFlag) String() string {
	return of.value.String()
}

func (of *OCIImageRefFlag) Type() string {
	return "oci-image"
}

var _ pflag.Value = &OCIImageRefFlag{}

func OCIImageRefVar(fs *pflag.FlagSet, ptr *meta.OCIImageRef, name string, defVal meta.OCIImageRef, usage string) {
	OCIImageRefVarP(fs, ptr, name, "", defVal, usage)
}

func OCIImageRefVarP(fs *pflag.FlagSet, ptr *meta.OCIImageRef, name, shorthand string, defVal meta.OCIImageRef, usage string) {
	fs.VarP(&OCIImageRefFlag{value: defVal}, name, shorthand, usage)
}
