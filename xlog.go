package xlog

import (
	"log"

	"github.com/sirupsen/logrus"
)

// LogLevel, is compatible with logrus
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel = iota
	// FatalLevel level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
)

// map to the logrus.Formatter
const (
	TextFormatter = iota
	JSONFormatter
)

// Logger is the abstract of logrus.FieldLogger and logrus.StdLogger
type Logger interface {
	logrus.FieldLogger
}

func levelToLogrusLevel(level int) logrus.Level {
	switch level {
	case PanicLevel:
		return logrus.PanicLevel
	case FatalLevel:
		return logrus.FatalLevel
	case ErrorLevel:
		return logrus.ErrorLevel
	case WarnLevel:
		return logrus.WarnLevel
	case InfoLevel:
		return logrus.InfoLevel
	case DebugLevel:
		return logrus.DebugLevel

	default:
		log.Fatalf("xlog: Invalid Loglevel: %d", level)
	}

	return logrus.InfoLevel
}

func createLogrusFormatter(formatter int) logrus.Formatter {
	switch formatter {
	case TextFormatter:
		return new(logrus.TextFormatter)
	case JSONFormatter:
		return new(logrus.JSONFormatter)
	default:
		log.Fatalf("xlog: Invalid LogFormatter: %d", formatter)
	}

	return nil
}
