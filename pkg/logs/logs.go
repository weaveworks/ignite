package logs

import (
	"io/ioutil"
	golog "log"
	"os"

	log "github.com/sirupsen/logrus"
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
