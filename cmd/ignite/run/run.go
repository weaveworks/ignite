package run

import (
	flag "github.com/spf13/pflag"
)

type RunFlags struct {
	*CreateFlags
	*StartFlags
}

type RunOptions struct {
	*CreateOptions
	*StartOptions
}

func (rf *RunFlags) NewRunOptions(args []string, fs *flag.FlagSet) (*RunOptions, error) {
	co, err := rf.NewCreateOptions(args, fs)
	if err != nil {
		return nil, err
	}

	// TODO: We should be able to use the constructor here instead...
	so := &StartOptions{
		StartFlags: rf.StartFlags,
		AttachOptions: &AttachOptions{
			checkRunning: false,
		},
	}

	return &RunOptions{co, so}, nil
}

func Run(ro *RunOptions, fs *flag.FlagSet) error {
	if err := Create(ro.CreateOptions); err != nil {
		return err
	}

	// Copy the pointer over for Start
	// TODO: This is pretty bad, fix this
	ro.vm = ro.VM

	return Start(ro.StartOptions, fs)
}
