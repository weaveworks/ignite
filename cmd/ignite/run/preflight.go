package run

import (
	"github.com/weaveworks/ignite/pkg/preflight/checkers"
	"github.com/weaveworks/ignite/pkg/util"
	"k8s.io/apimachinery/pkg/util/sets"
)

type PreflightFlags struct {
	IgnoredPreflightErrors []string
}

type preflightOptions struct {
	*PreflightFlags
}

func (pf *PreflightFlags) NewPreflightOptions() preflightOptions {
	return preflightOptions{
		PreflightFlags: pf,
	}
}

func Preflight(po preflightOptions) error {
	ignoredPreflightErrors := sets.NewString(util.ToLower(po.IgnoredPreflightErrors)...)
	if err := checkers.PreflightCmdChecks(ignoredPreflightErrors); err != nil {
		return err
	}
	return nil
}
