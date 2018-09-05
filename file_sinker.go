package xlog

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

// The time based rotation strategy constant
const (
	TimeRotateNone = iota
	TimeRotateHourly
	TimeRotateDaily
)

// FileSinkerConf holds the configurable attributes of FileSinker
type FileSinkerConf struct {
	// Whether flush the log data to disk as soon as it arrived
	// Default false
	FlushImmediately bool

	// Default 0, diable rotation by size
	RotateSize int

	// Default TimeRotateNone, diable rotation by time
	RotateTime int

	// The directory of the log data to sink.
	// Note: the last part of the path will be treated as the prefix of the log files
	//       if the Path ends with a "/", the "log" string will be used as the prefix
	Path string
}

// NewFileSinkerConf creates a FileSinkerConf instance with default value
func NewFileSinkerConf(p string) *FileSinkerConf {
	if len(p) <= 0 {
		return nil
	}

	return &FileSinkerConf{
		FlushImmediately: false,
		RotateSize:       0,
		RotateTime:       TimeRotateNone,
		Path:             p,
	}
}

// FileSinker implements the logrus.Logger.Out interface.
// It has the following feature:
// 1. Rotate log files according time
// 2. Rotate log files according file size
//
// The "Sinker" concept is from apache.flume
type FileSinker struct {
	conf *FileSinkerConf

	outer *bufio.Writer // outer is the handle of the current log file

	dir    string
	prefix string

	sizeWritten int      // the size of the current file
	curfile     string   // the current file name
	curfp       *os.File // the file-pointer of the current file

	fseq int    // the current file sequence number while rotate by size
	tstr string // the current time string while rotate by time

	// ch is the log data channel
	// ch   chan *[]byte
	// quit chan struct{}
}

// NewFileSinker creates a FileSinker instance with the given
// conf. It returns nil if the conf is invalid or some other error happend
func NewFileSinker(conf *FileSinkerConf) *FileSinker {
	if nil == conf {
		return nil
	}

	dir, prefix := path.Split(conf.Path)
	if len(dir) <= 0 {
		log.Printf("FileSinker path is empty")
		return nil
	}

	if _, err := os.Stat(dir); err != nil && !os.IsNotExist(err) {
		log.Printf("FileSinker path failer: [%s: %s]", dir, err.Error())
		return nil
	} else if err != nil {
		if err1 := os.MkdirAll(dir, 0755); err1 != nil {
			log.Printf("FileSinker mkdir failed: [%s: %s]", dir, err.Error())
			return nil
		}
	}

	if len(prefix) <= 0 {
		prefix = "log"
	}

	fsk := &FileSinker{
		conf:        conf,
		outer:       nil,
		dir:         dir[:len(dir)-1], // we drop the last "/" in dirpath
		prefix:      prefix,
		sizeWritten: -1,
		curfp:       nil,
		curfile:     "",
		fseq:        -1, // set it to -1 so that the size rotation could be triggered for the first time
		tstr:        "",
		//	ch:     make(chan *[]byte, 1024),
		//	quit:   make(chan struct{}),
	}

	return fsk
}

// Close flush the buffer and close the file
func (fsk *FileSinker) Close() {
	if fsk.outer != nil {
		fsk.outer.Flush()
		fsk.outer = nil
	}

	if fsk.curfp != nil {
		fsk.curfp.Close()
		fsk.curfp = nil
	}

	fsk.sizeWritten = -1
	fsk.curfile = ""
	fsk.fseq = -1
	fsk.tstr = ""
}

// Write implements the logrus.Logger.Out interface
func (fsk *FileSinker) Write(p []byte) (n int, err error) {
	if p == nil {
		return
	}

	if fsk.outer == nil {
		fsk.rotate()
	}

	if fsk.outer == nil {
		log.Panicf("FileSinker create writer failed")
	}

	n, err = fsk.outer.Write(p)
	if err != nil {
		fsk.sizeWritten += n
		if fsk.conf.FlushImmediately {
			fsk.outer.Flush()
		}
	}

	fsk.rotate()

	return n, err
}

// The rotate function is invoked only in the Write function.
// All the Wirte operation will be protected by the logrus.Logger's lock,
// Therefore, we need not use another lock to protect the fsk.outer and
// fsk.curfp
func (fsk *FileSinker) rotate() {
	// if no rotation is configured, the log file name will be:
	//		prefix
	// if rotate by size, the file name will be:
	//		prefix.0, prefix.1, ...
	// if rotate by time, the file name will be:
	//		prefix.YYYY-MM-DD[-HH]
	// if rotate by both time and size, the file name will be:
	//		prefix.YYYY-MM-DD[-HH].0, prefix.YYYY-MM-DD[-HH].1, ...

	filename := fmt.Sprintf("%s/%s", fsk.dir, fsk.prefix)

	if fsk.conf.RotateSize <= 0 &&
		fsk.conf.RotateTime == TimeRotateNone {
		fsk.fileRotate(filename)
		return
	}

	var tstr string
	if fsk.conf.RotateTime != TimeRotateNone {
		switch fsk.conf.RotateTime {
		case TimeRotateHourly:
			tstr = time.Now().Format("2006-01-02-15")
		case TimeRotateDaily:
			tstr = time.Now().Format("2006-01-02")
		default:
		}

		filename = fmt.Sprintf("%s.%s", filename, tstr)
	}

	fseq := fsk.fseq
	if fsk.conf.RotateSize > 0 {
		if tstr != fsk.tstr {
			fseq = 0
		} else if fsk.sizeWritten == -1 ||
			fsk.sizeWritten >= fsk.conf.RotateSize {
			fseq++
		}

		fseq = fsk.findAvailableFileSeq(fseq, filename)
		filename = fmt.Sprintf("%s.%d", filename, fseq)
	}

	if fsk.fileRotate(filename) {
		fsk.fseq = fseq
		fsk.tstr = tstr
	}
}

func (fsk *FileSinker) fileRotate(filename string) bool {
	if filename != fsk.curfile {
		fp, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("FileSinker open new file failed: [%s: %s]", filename, err.Error())
			return false
		}

		if fsk.outer != nil {
			fsk.outer.Flush()
		}

		if fsk.curfp != nil {
			fsk.curfp.Close()
		}

		fsk.sizeWritten = 0
		fsk.curfile = filename
		fsk.outer = bufio.NewWriter(fp)
	}

	return true
}

// Find the first available file sequence in case that the previous sequence numbers
// have already been used
func (fsk *FileSinker) findAvailableFileSeq(fseq int, filename string) int {
	if fseq != fsk.fseq {
		for seq := fseq; seq < 1000000; seq++ {
			file := fmt.Sprintf("%s.%d", filename, seq)
			if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
				return seq
			}
		}

		log.Panicf("FileSinker could not find an available sequence number")
	}

	return fseq
}
