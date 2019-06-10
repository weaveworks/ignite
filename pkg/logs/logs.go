package logs

import (
	"fmt"
	"io/ioutil"
	golog "log"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// Quiet specifies whether to only print machine-readable IDs
var Quiet bool

// InitLogs initializes the logging system for ignite
func InitLogs() {
	// Initialize a new logrus logger
	var Logger = log.New()

	// Set the output to be stdout in the normal case, but discard all log output in quiet mode
	Logger.SetOutput(os.Stdout)
	if Quiet {
		Logger.SetOutput(ioutil.Discard)
	}

	// Disable timestamp logging, but still output the seconds elapsed
	Logger.SetFormatter(&log.TextFormatter{
		DisableTimestamp: false,
		FullTimestamp:    false,
	})

	// Disable the stdlib's automatic add of the timestamp in beginning of the log message,
	// as we stream the logs from stdlib log to this logrus instance.
	golog.SetFlags(0)
	golog.SetOutput(Logger.Writer())

	// Also forward all logs from the standard logrus logger to our specific instance
	log.StandardLogger().SetOutput(Logger.Writer())
}

// PrintMachineReadableID prints the machine-readable ID if we're in quiet mode, otherwise is a no-op
func PrintMachineReadableID(id string, err error) error {
	if err != nil {
		return err
	}
	if Quiet {
		fmt.Println(id)
	}
	return nil
}

// AddQuietFlag adds the quiet flag to a flagset
func AddQuietFlag(fs *pflag.FlagSet) {
	fs.BoolVarP(&Quiet, "quiet", "q", Quiet, "The quiet mode allows for machine-parsable output, by printing only IDs")
}
