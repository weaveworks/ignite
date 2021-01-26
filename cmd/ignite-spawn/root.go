package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
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

	if util.IDPrefix == "" {
		util.IDPrefix = constants.IGNITE_PREFIX
	}

	util.GenericCheckErr(func() error {
		vm, err := decodeVM(fs.Args()[0])
		if err != nil {
			return err
		}

		return StartVM(vm)
	}())
}

func usage() {
	util.GenericCheckErr(fmt.Errorf("usage: ignite-spawn [--log-level <level>] [--id-prefix <prefix>] <vm>"))
}

func addGlobalFlags(fs *pflag.FlagSet) {
	// TODO: Add a version flag
	logflag.LogLevelFlagVar(fs, &logLevel)
	cmdutil.AddIDPrefixFlag(fs, &util.IDPrefix)
}
