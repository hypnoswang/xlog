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

// Close gives an chance to do some cleaning work
// You must invoke this function carefully, because we don't have any
// synchronous machanism to protect the sinker's outer and curfp.
// A good way is that close this logger just before your program exit.
func (f *FileLogger) Close() {
	sinker := f.Logger.Out.(*FileSinker)
	sinker.Close()
}
