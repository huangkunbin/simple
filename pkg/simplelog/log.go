package simplelog

import (
	"fmt"
	"reflect"
	"sync"
)

type LoggerLevel int

const (
	DebugLv LoggerLevel = 0
	InfoLv  LoggerLevel = 1
	WarnLv  LoggerLevel = 2
	ErrorLv LoggerLevel = 3
)

type loggerTypeMap map[reflect.Type]Logger

var (
	loggerStorage = make(loggerTypeMap)
	mu            sync.Mutex
	loggerLv      = DebugLv
	DefaultLogger = NewConsoleLogger()
)

func SetLogLevel(lv LoggerLevel) {
	if lv != loggerLv {
		Infof("更改日志级别：%d -> %d", loggerLv, lv)
		loggerLv = lv
	}
}

// RegisterLogger 注册日志记录器;支持重复注册
func RegisterLogger(loggers ...Logger) {
	mu.Lock()
	defer mu.Unlock()
	for _, l := range loggers {
		logTyp := reflect.TypeOf(l)
		if logTyp.Kind() == reflect.Ptr {
			logTyp = logTyp.Elem()
		}
		if _, has := loggerStorage[logTyp]; !has {
			loggerStorage[logTyp] = l
			Infof("注册logger[%s]\n", logTyp)
		}
	}
}

func UnRegisterLogger(loggers ...Logger) {
	mu.Lock()
	defer mu.Unlock()
	for _, l := range loggers {
		logTyp := reflect.TypeOf(l)
		if logTyp.Kind() == reflect.Ptr {
			logTyp = logTyp.Elem()
		}
		if _, has := loggerStorage[logTyp]; has {
			delete(loggerStorage, logTyp)
			fmt.Printf("取消logger[%s]\n", logTyp)
		}
	}
}

func (lm loggerTypeMap) call(lv LoggerLevel, fn func(Logger)) {
	if loggerLv > lv {
		return
	}
	for _, logger := range lm {
		fn(logger)
	}
}

func Debug(v ...interface{}) {
	loggerStorage.call(DebugLv, func(l Logger) {
		l.Debug(v...)
	})
}

func Debugf(format string, v ...interface{}) {
	loggerStorage.call(DebugLv, func(l Logger) {
		l.Debugf(format, v...)
	})
}

func Info(v ...interface{}) {
	loggerStorage.call(InfoLv, func(l Logger) {
		l.Info(v...)
	})
}

func Infof(format string, v ...interface{}) {
	loggerStorage.call(InfoLv, func(l Logger) {
		l.Infof(format, v...)
	})
}

func Warn(v ...interface{}) {
	loggerStorage.call(WarnLv, func(l Logger) {
		l.Warn(v...)
	})
}

func Warnf(format string, v ...interface{}) {
	loggerStorage.call(WarnLv, func(l Logger) {
		l.Warnf(format, v...)
	})
}

func Error(v ...interface{}) {
	loggerStorage.call(ErrorLv, func(l Logger) {
		l.Error(v...)
	})
}

func Errorf(format string, v ...interface{}) {
	loggerStorage.call(ErrorLv, func(l Logger) {
		l.Errorf(format, v...)
	})
}
