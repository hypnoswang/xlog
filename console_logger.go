package xlog

import (
	"os"

	"github.com/sirupsen/logrus"
)

// ConsoleLogger put log contents to stdout
type ConsoleLogger struct {
	logrus.Logger
}

// NewConsoleLogger create a new ConsoleLogger
// TODO: add hooks parameters
func NewConsoleLogger(level, formatter int) *ConsoleLogger {
	lvl := levelToLogrusLevel(level)
	fmtr := createLogrusFormatter(formatter)

	logger := &ConsoleLogger{
		Logger: logrus.Logger{
			Out:       os.Stdout,
			Formatter: fmtr,
			Hooks:     make(logrus.LevelHooks),
			Level:     lvl,
		},
	}

	return logger
}

func (c *ConsoleLogger) Close() {
	// do nothing
}
