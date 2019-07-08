package run

import (
	"fmt"

	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/filter"
	"github.com/weaveworks/ignite/pkg/storage/filterer"
)

type RunFlags struct {
	*CreateFlags
	*StartFlags
}

type runOptions struct {
	*createOptions
	*startOptions
}

func (rf *RunFlags) NewRunOptions(args []string) (*runOptions, error) {
	// parse the args and the config file
	err := rf.CreateFlags.parseArgsAndConfig(args)
	if err != nil {
		return nil, err
	}

	imageName := rf.VM.Spec.Image.OCIClaim.Ref

	// Logic to import the image if it doesn't exist
	if _, err := client.Images().Find(filter.NewIDNameFilter(imageName)); err != nil { // TODO: Use this match in create?
		switch err.(type) {
		case *filterer.NonexistentError:
			io, err := NewImportOptions(imageName)
			if err != nil {
				return nil, err
			}

			if err := Import(io); err != nil {
				return nil, err
			}
		default:
			return nil, err
		}
	}

	co, err := rf.NewCreateOptions(args)
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
	fmt.Println("create")
	if err := Create(ro.createOptions); err != nil {
		return err
	}

	// Copy the pointer over for Start
	// TODO: This is pretty bad, fix this
	ro.vm = ro.newVM

	fmt.Println("start")
	return Start(ro.startOptions)
}
