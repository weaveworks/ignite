package cmdutil

import (
	"github.com/spf13/pflag"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

type SizeFlag struct {
	value *meta.Size
}

func (sf *SizeFlag) Set(val string) error {
	var err error
	*sf.value, err = meta.NewSizeFromString(val)
	return err
}

func (sf *SizeFlag) String() string {
	return sf.value.String()
}

func (sf *SizeFlag) Type() string {
	return "size"
}

var _ pflag.Value = &SizeFlag{}

func SizeVar(fs *pflag.FlagSet, ptr *meta.Size, name, usage string) {
	SizeVarP(fs, ptr, name, "", usage)
}

func SizeVarP(fs *pflag.FlagSet, ptr *meta.Size, name, shorthand, usage string) {
	fs.VarP(&SizeFlag{value: ptr}, name, shorthand, usage)
}
