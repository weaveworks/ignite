package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/pkg/logs"
)

var logLevel = logrus.InfoLevel

// RunIgniteSpawn runs the root command for ignite-spawn
func RunIgniteSpawn() {
	fs := &pflag.FlagSet{
		Usage: usage,
	}

	addGlobalFlags(fs)
	checkErr(fs.Parse(os.Args[1:]))
	logs.Logger.SetLevel(logLevel)

	if len(fs.Args()) != 1 {
		usage()
	}

	checkErr(func() error {
		vm, err := decodeVM(fs.Args()[0])
		if err != nil {
			return err
		}

		return StartVM(vm)
	}())
}

func usage() {
	checkErr(fmt.Errorf("usage: ignite-spawn [--log-level <level>] <vm>"))
}

func addGlobalFlags(fs *pflag.FlagSet) {
	cmdutil.LogLevelFlagVar(fs, &logLevel)
}

// checkErr is used by the ignite-spawn command to check if the action failed
// and respond with a fatal error provided by the logger (calls os.Exit)

func checkErr(err error) {
	switch err.(type) {
	case nil:
		return // Don't fail if there's no error
	}

	logrus.Fatal(err)
}
