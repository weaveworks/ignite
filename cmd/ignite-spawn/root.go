package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/pkg/logs"
)

var logLevel = logrus.InfoLevel

// RunIgniteSpawn runs the root command for ignite-spawn
func RunIgniteSpawn() {
	fs := &pflag.FlagSet{
		Usage: usage,
	}

	addGlobalFlags(fs)
	cmdutil.CheckErr(fs.Parse(os.Args[1:]))
	logs.Logger.SetLevel(logLevel)

	if len(fs.Args()) != 1 {
		usage()
	}

	cmdutil.CheckErr(func() error {
		opts, err := NewOptions(fs.Args()[0])
		if err != nil {
			return err
		}

		return StartVM(opts)
	}())
}

func usage() {
	cmdutil.CheckErr(fmt.Errorf("usage: ignite-spawn [--log-level <level>] <vm>"))
}

func addGlobalFlags(fs *pflag.FlagSet) {
	cmdutil.LogLevelFlagVar(fs, &logLevel)
}
