package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/logs"
	logflag "github.com/weaveworks/ignite/pkg/logs/flag"
	"github.com/weaveworks/ignite/pkg/util"
)

var logLevel = logrus.InfoLevel

// RunIgniteSpawn runs the root command for ignite-spawn
func RunIgniteSpawn() {
	fs := &pflag.FlagSet{
		Usage: usage,
	}

	addGlobalFlags(fs)
	util.GenericCheckErr(fs.Parse(os.Args[1:]))
	logs.Logger.SetLevel(logLevel)

	if len(fs.Args()) != 1 {
		usage()
	}

	// there is no default value so this should be initialized if the VM has no status prefix already (backwards-compat)
	if util.IDPrefix == "" {
		util.IDPrefix = constants.IGNITE_PREFIX
	}

	util.GenericCheckErr(func() error {
		vm, err := decodeVM(fs.Args()[0])
		if err != nil {
			return err
		}

		// guard against no status -- otherwise override default
		if vm.Status.IDPrefix != "" {
			util.IDPrefix = vm.Status.IDPrefix
		}

		return StartVM(vm)
	}())
}

func usage() {
	util.GenericCheckErr(fmt.Errorf("usage: ignite-spawn [--log-level <level>] <vm>"))
}

func addGlobalFlags(fs *pflag.FlagSet) {
	// TODO: Add a version flag
	logflag.LogLevelFlagVar(fs, &logLevel)
}
