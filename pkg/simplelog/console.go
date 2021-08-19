package simplelog

import (
	"fmt"
	"log"
	"os"
)

const calldepth = 3

type consoleLogger struct {
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	err   *log.Logger
}

func NewConsoleLogger() Logger {
	return &consoleLogger{
		debug: log.New(os.Stdout, "\033[40;32m[DEBUG] ", log.LstdFlags|log.Lshortfile),
		info:  log.New(os.Stdout, "\033[40;36m[INFO ] ", log.LstdFlags|log.Lshortfile),
		warn:  log.New(os.Stdout, "\033[40;33m[WARN ] ", log.LstdFlags|log.Lshortfile),
		err:   log.New(os.Stdout, "\033[41;37m[ERROR] ", log.LstdFlags|log.Lshortfile),
	}
}

func (l *consoleLogger) Debug(v ...interface{}) {
	l.debug.Output(calldepth, "\033[0m"+fmt.Sprint(v...))
}

func (l *consoleLogger) Debugf(format string, v ...interface{}) {
	l.debug.Output(calldepth, "\033[0m"+fmt.Sprintf(format, v...))
}

func (l *consoleLogger) Info(v ...interface{}) {
	l.info.Output(calldepth, "\033[0m"+fmt.Sprint(v...))
}

func (l *consoleLogger) Infof(format string, v ...interface{}) {
	l.info.Output(calldepth, "\033[0m"+fmt.Sprintf(format, v...))
}

func (l *consoleLogger) Warn(v ...interface{}) {
	l.warn.Output(calldepth, "\033[0m"+fmt.Sprint(v...))
}

func (l *consoleLogger) Warnf(format string, v ...interface{}) {
	l.warn.Output(calldepth, "\033[0m"+fmt.Sprintf(format, v...))
}

func (l *consoleLogger) Error(v ...interface{}) {
	l.err.Output(calldepth, "\033[0m"+fmt.Sprint(v...))
}

func (l *consoleLogger) Errorf(format string, v ...interface{}) {
	l.err.Output(calldepth, "\033[0m"+fmt.Sprintf(format, v...))
}
