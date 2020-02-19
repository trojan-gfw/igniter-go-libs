package simpleandroidlog

import (
	"io"
	golog "log"

	"github.com/eycorsican/go-tun2socks/common/log"
)

var mylogger AndroidLogger

type AndroidLogger interface {
	GetLevel() log.LogLevel
	log.Logger
}

type simpleAndroidLogger struct {
	level  log.LogLevel
	prefix string
}

func NewSimpleAndroidLogger() AndroidLogger {
	return &simpleAndroidLogger{
		level:  log.INFO,
		prefix: "[tun2socks] ",
	}
}

func (l *simpleAndroidLogger) SetLevel(level log.LogLevel) {
	l.level = level
}

func (l *simpleAndroidLogger) GetLevel() log.LogLevel {
	return l.level
}

func (l *simpleAndroidLogger) Debugf(msg string, args ...interface{}) {
	if l.level <= log.DEBUG {
		l.output(msg, args...)
	}
}

func (l *simpleAndroidLogger) Infof(msg string, args ...interface{}) {
	if l.level <= log.INFO {
		l.output(msg, args...)
	}
}

func (l *simpleAndroidLogger) Warnf(msg string, args ...interface{}) {
	if l.level <= log.WARN {
		l.output(msg, args...)
	}
}

func (l *simpleAndroidLogger) Errorf(msg string, args ...interface{}) {
	if l.level <= log.ERROR {
		l.output(msg, args...)
	}
}

func (l *simpleAndroidLogger) Fatalf(msg string, args ...interface{}) {
	golog.Fatalf(l.prefix+msg, args...)
}

func (l *simpleAndroidLogger) output(msg string, args ...interface{}) {
	golog.Printf(l.prefix+msg, args...)
}

func (l *simpleAndroidLogger) GetUnderlyingWriter() io.Writer {
	return golog.Writer()
}

func GetLogger() AndroidLogger {
	return mylogger
}

func init() {
	mylogger = NewSimpleAndroidLogger()
	log.RegisterLogger(mylogger)
}
