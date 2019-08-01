package flag

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

type LogLevelFlag struct {
	value *logrus.Level
}

func (lf *LogLevelFlag) Set(val string) error {
	var err error
	*lf.value, err = logrus.ParseLevel(val)
	return err
}

func (lf *LogLevelFlag) String() string {
	if lf.value == nil {
		return ""
	}
	return lf.value.String()
}

func (lf *LogLevelFlag) Type() string {
	return "loglevel"
}

var _ pflag.Value = &LogLevelFlag{}

func LogLevelFlagVar(fs *pflag.FlagSet, ptr *logrus.Level) {
	fs.Var(&LogLevelFlag{value: ptr}, "log-level", "Specify the loglevel for the program")
}
