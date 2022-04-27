package logs

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	z *zap.Logger
)

func GlobalLogger() *zap.Logger {
	return z
}

func InitLogger(caller int, opts ...OutputOption) *zap.Logger {
	if len(opts) == 0 {
		opts = append(opts, NewConsoleOption(zap.DebugLevel), DefaultFile())
	}
	z = NewZapLog(caller, opts...)
	return z
}

func init() {
	InitLogger(1)
}

type OutputOption struct {
	Encoder zapcore.Encoder
	Writer  zapcore.WriteSyncer
	Level   zapcore.Level
}

func NewZapLog(caller int, outputs ...OutputOption) *zap.Logger {
	var cores []zapcore.Core
	for _, v := range outputs {
		core := zapcore.NewCore(v.Encoder, v.Writer, v.Level)
		cores = append(cores, core)
	}
	logger := zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(caller))
	return logger
}

func Debug(f interface{}, v ...interface{}) {
	z.Debug(FormatLog(f, v...))
}

func Warn(f interface{}, v ...interface{}) {
	z.Warn(FormatLog(f, v...))
}

func Info(f interface{}, v ...interface{}) {
	z.Info(FormatLog(f, v...))
}

func Error(f interface{}, v ...interface{}) {
	z.Error(FormatLog(f, v...))
}

func FormatLog(f interface{}, v ...interface{}) string {
	var msg string
	switch f.(type) {
	case string:
		msg = f.(string)
		if len(v) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
			//format string
		} else {
			//do not contain format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	msg = fmt.Sprintf(msg, v...)
	return msg
}
