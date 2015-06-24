// Package log contains functions for global debug and logging.
package log

import (
	"github.com/op/go-logging"
	"os"
)

const (
	CRITICAL int = 0
	ERROR        = 1
	WARNING      = 2
	NOTICE       = 3
	INFO         = 4
	DEBUG        = 5
)

// set up the global logger
var log = logging.MustGetLogger("global")

var format = logging.MustStringFormatter(
	"%{color} %{level:.4s} ▶ %{color:reset} %{message}",
)

var backend = logging.NewLogBackend(os.Stderr, "", 0)
var backendFormatted = logging.NewBackendFormatter(backend, format)
var backendLeveled = logging.AddModuleLevel(backendFormatted)

func init() {
	logging.SetBackend(backendLeveled)
	SetLogLevel(NOTICE)
}

// Log logs a formatted message with the specified integer verbosity level.
// The lower the level, the more critical the message.
func Log(level int, format string, args ...interface{}) {
	switch level {
	case CRITICAL:
		log.Critical(format, args...)
	case ERROR:
		log.Error(format, args...)
	case WARNING:
		log.Warning(format, args...)
	case NOTICE:
		log.Notice(format, args...)
	case INFO:
		log.Info(format, args...)
	case DEBUG:
		log.Debug(format, args...)
	default:
		log.Error(format, args...)
	}
}

// SetLogLevel sets the verbosity level of the logger, with 0 being least verbose,
// and 5 being the most verbose. By default, the verbosity level is 1, which
// logs critical and error messages.
func SetLogLevel(level int) {
	switch level {
	case CRITICAL:
		backendLeveled.SetLevel(logging.CRITICAL, "global")
	case ERROR:
		backendLeveled.SetLevel(logging.ERROR, "global")
	case WARNING:
		backendLeveled.SetLevel(logging.WARNING, "global")
	case NOTICE:
		backendLeveled.SetLevel(logging.NOTICE, "global")
	case INFO:
		backendLeveled.SetLevel(logging.INFO, "global")
	case DEBUG:
		backendLeveled.SetLevel(logging.DEBUG, "global")
	default:
		backendLeveled.SetLevel(logging.ERROR, "global")
	}

}
