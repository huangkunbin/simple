package simplelog

import (
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *zap.SugaredLogger
}

func NewFileLogger(savepath, prefix string) Logger {
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02T15:04:05.000"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})
	core := zapcore.NewCore(encoder, zapcore.AddSync(func() io.Writer {
		hook, err := rotatelogs.New(path.Join(savepath, prefix+"-%Y%m%d%H.log"))
		if err != nil {
			panic(err)
		}
		return hook
	}()), zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	}))
	logger := zap.New(core).Sugar()
	return &zapLogger{logger: logger}
}

func (l *zapLogger) Debug(v ...interface{}) {
	l.logger.Debug(l.noWrap(v...))
}

func (l *zapLogger) Debugf(format string, v ...interface{}) {
	l.logger.Debug(l.noWrapf(format, v...))
}

func (l *zapLogger) Info(v ...interface{}) {
	l.logger.Info(l.noWrap(v...))
}

func (l *zapLogger) Infof(format string, v ...interface{}) {
	l.logger.Info(l.noWrapf(format, v...))
}

func (l *zapLogger) Warn(v ...interface{}) {
	l.logger.Warn(l.noWrap(v...))
}

func (l *zapLogger) Warnf(format string, v ...interface{}) {
	l.logger.Warn(l.noWrapf(format, v...))
}

func (l *zapLogger) Error(v ...interface{}) {
	l.logger.Error(l.noWrap(v...))
}

func (l *zapLogger) Errorf(format string, v ...interface{}) {
	l.logger.Error(l.noWrapf(format, v...))
}

func (l *zapLogger) noWrap(v ...interface{}) string {
	return l.newLine(strings.ReplaceAll(fmt.Sprint(v...), "", ""))
}

func (l *zapLogger) noWrapf(format string, v ...interface{}) string {
	return l.newLine(strings.ReplaceAll(fmt.Sprintf(format, v...), "", ""))
}

func (l *zapLogger) newLine(s string) string {
	return strings.TrimRight(s, "\n")
}
