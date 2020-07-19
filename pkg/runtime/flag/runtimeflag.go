package flag

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/pkg/runtime"
)

var runtimes = runtime.ListRuntimes()

type RuntimeFlag struct {
	value *runtime.Name
}

func (nf *RuntimeFlag) Set(val string) error {
	for _, r := range runtimes {
		if r.String() == val {
			*nf.value = r
			return nil
		}
	}

	return fmt.Errorf("invalid runtime %q, must be one of %v", val, runtimes)
}

func (nf *RuntimeFlag) String() string {
	if nf.value == nil {
		return ""
	}

	return nf.value.String()
}

func (nf *RuntimeFlag) Type() string {
	return "runtime"
}

var _ pflag.Value = &RuntimeFlag{}

func RuntimeVar(fs *pflag.FlagSet, ptr *runtime.Name) {
	fs.Var(&RuntimeFlag{value: ptr}, "runtime", fmt.Sprintf("Container runtime to use. Available options are: %v (default %v)", runtimes, runtime.RuntimeContainerd))
}
