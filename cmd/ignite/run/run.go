package run

import (
	"github.com/weaveworks/ignite/cmd/ignite/run/runutil"
	"github.com/weaveworks/ignite/pkg/metadata"
)

type RunFlags struct {
	*CreateFlags
	*StartFlags
}

type runOptions struct {
	*createOptions
	*startOptions
}

func (rf *RunFlags) NewRunOptions(l *runutil.ResLoader, imageMatch string) (*runOptions, error) {
	// Logic to import the image if it doesn't exist
	if allImages, err := l.Images(); err == nil {
		if _, err := allImages.MatchSingle(imageMatch); err != nil { // TODO: Use this match in create?
			if _, ok := err.(*metadata.NonexistentError); !ok {
				return nil, err
			}

			io, err := (&ImportFlags{}).NewImportOptions(l, imageMatch)
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

	co, err := rf.NewCreateOptions(l, imageMatch)
	if err != nil {
		return nil, err
	}

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
	ro.vm = ro.newVM

	if err := Start(ro.startOptions); err != nil {
		return err
	}

	return nil
}
