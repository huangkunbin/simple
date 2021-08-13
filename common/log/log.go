package log

import (
	"os"
	"path/filepath"
	"simple/lib/simplelog"
	"strings"
)

var (
	Debug  = simplelog.Debug
	Debugf = simplelog.Debugf
	Info   = simplelog.Info
	Infof  = simplelog.Infof
	Warn   = simplelog.Warn
	Warnf  = simplelog.Warnf
	Error  = simplelog.Error
	Errorf = simplelog.Errorf
)

// 输出平台枚举: 二进制按位与进行标记
const (
	ConsoleLog = 1 << 0 // 控制台输出
	FileLog    = 1 << 1 // 文件输出
	SLSLog     = 1 << 2 // 阿里云日志服务
)

var (
	slsLogger  simplelog.SLSLogger
	fileLogger simplelog.Logger
)

func init() {
	simplelog.SetLogLevel(simplelog.DebugLv)
	simplelog.RegisterLogger(simplelog.DefaultLogger)
	fileLogger = NewFileLogger()
}

func NewFileLogger() simplelog.Logger {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	prefix, err := os.Executable()
	if err != nil {
		panic(err)
	}
	prefix = filepath.Base(prefix)
	prefix = strings.Replace(prefix, "\\", "/", -1)
	return simplelog.NewFileLogger(dir+"/logs", prefix)
}

func UpdateLoggers(loggerLevel, loggerType int) {
	simplelog.SetLogLevel(simplelog.LoggerLevel(loggerLevel))
	// 控制台
	if (loggerType & ConsoleLog) != 0 {
		simplelog.RegisterLogger(simplelog.DefaultLogger)
	} else {
		simplelog.UnRegisterLogger(simplelog.DefaultLogger)
	}
	// 文件
	if (loggerType & FileLog) != 0 {
		simplelog.RegisterLogger(fileLogger)
	} else {
		simplelog.UnRegisterLogger(fileLogger)
	}
}
