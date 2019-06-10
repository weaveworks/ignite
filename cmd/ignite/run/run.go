package run

import "github.com/weaveworks/ignite/cmd/ignite/run/runutil"

type RunFlags struct {
	*CreateFlags
	*StartFlags
}

type runOptions struct {
	*createOptions
	*startOptions
}

func (rf *RunFlags) NewRunOptions(l *runutil.ResourceLoader, imageMatch string) (*runOptions, error) {
	co, err := rf.NewCreateOptions(l, imageMatch)
	if err != nil {
		return nil, err
	}

	so := &startOptions{
		StartFlags:    rf.StartFlags,
		attachOptions: &attachOptions{
			checkRunning: false,
			vm:           co.newVM, // This should hopefully work
		},
	}

	return &runOptions{co, so}, nil
}

func Run(ro *runOptions) error {
	if err := Create(ro.createOptions); err != nil {
		return err
	}

	if err := Start(ro.startOptions); err != nil {
		return err
	}

	return nil
}