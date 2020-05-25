package run

import (
	flag "github.com/spf13/pflag"
)

type RunFlags struct {
	*CreateFlags
	*StartFlags
}

type runOptions struct {
	*createOptions
	*startOptions
}

func (rf *RunFlags) NewRunOptions(args []string, fs *flag.FlagSet) (*runOptions, error) {
	co, err := rf.NewCreateOptions(args, fs)
	if err != nil {
		return nil, err
	}

	// TODO: We should be able to use the constructor here instead...
	so := &startOptions{
		StartFlags: rf.StartFlags,
		attachOptions: &attachOptions{
			checkRunning: false,
		},
	}

	return &runOptions{co, so}, nil
}

func Run(ro *runOptions) error {
	if err := Create(ro.createOptions); err != nil {
		return err
	}

	// Copy the pointer over for Start
	// TODO: This is pretty bad, fix this
	ro.vm = ro.VM

	return Start(ro.startOptions)
}
