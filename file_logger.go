package xlog

import (
	"github.com/sirupsen/logrus"
)

// FileLogger implements the function that put the log content to
// files with file-rotation
type FileLogger struct {
	logrus.Logger
}

// NewFileLogger create a new FileLogger
// TODO: add hooks support
func NewFileLogger(conf *FileSinkerConf, level, formatter int) *FileLogger {

	lvl := levelToLogrusLevel(level)
	fmtr := createLogrusFormatter(formatter)

	sinker := NewFileSinker(conf)
	if sinker == nil {
		return nil
	}

	logger := &FileLogger{
		Logger: logrus.Logger{
			Out:       sinker,
			Formatter: fmtr,
			Hooks:     make(logrus.LevelHooks),
			Level:     lvl,
		},
	}

	return logger
}
