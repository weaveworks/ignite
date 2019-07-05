package run

import (
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/loader"
)

type RunFlags struct {
	*CreateFlags
	*StartFlags
}

type runOptions struct {
	*createOptions
	*startOptions
}

func (rf *RunFlags) NewRunOptions(l *loader.ResLoader, args []string) (*runOptions, error) {
	// parse the args and the config file
	err := rf.CreateFlags.parseArgsAndConfig(args)
	if err != nil {
		return nil, err
	}

	// Logic to import the image if it doesn't exist
	if allImages, err := l.Images(); err == nil {
		imageName := rf.VM.Spec.Image.Ref
		if _, err := allImages.MatchSingle(imageName); err != nil { // TODO: Use this match in create?
			if _, ok := err.(*metadata.NonexistentError); !ok {
				return nil, err
			}

			io, err := NewImportOptions(l, imageName)
			if err != nil {
				return nil, err
			}

			if err := Import(io); err != nil {
				return nil, err
			}
		}
	} else {
		return nil, err
	}

	co, err := rf.NewCreateOptions(l, args)
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
	ro.vm = ro.newVM

	return Start(ro.startOptions)
}
