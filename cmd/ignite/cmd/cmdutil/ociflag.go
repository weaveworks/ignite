package cmdutil

import (
	"github.com/spf13/pflag"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

type OCIImageRefFlag struct {
	value *meta.OCIImageRef
}

func (of *OCIImageRefFlag) Set(val string) error {
	var err error
	*of.value, err = meta.NewOCIImageRef(val)
	return err
}

func (of *OCIImageRefFlag) String() string {
	if of.value == nil {
		return ""
	}
	return of.value.String()
}

func (of *OCIImageRefFlag) Type() string {
	return "oci-image"
}

var _ pflag.Value = &OCIImageRefFlag{}

func OCIImageRefVar(fs *pflag.FlagSet, ptr *meta.OCIImageRef, name, usage string) {
	OCIImageRefVarP(fs, ptr, name, "", usage)
}

func OCIImageRefVarP(fs *pflag.FlagSet, ptr *meta.OCIImageRef, name, shorthand, usage string) {
	fs.VarP(&OCIImageRefFlag{value: ptr}, name, shorthand, usage)
}
