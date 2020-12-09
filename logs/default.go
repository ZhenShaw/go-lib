package logs

import "go.uber.org/zap"

var l *zap.SugaredLogger

func init() {
	z := NewZapLog(0, DefaultConsole(), DefaultFile())
	l = z.Sugar()
}

func Debug(v ...interface{}) {
	l.Debug(v...)
}

func Warn(v ...interface{}) {
	l.Warn(v...)
}

func Info(v ...interface{}) {
	l.Info(v...)
}

func Error(v ...interface{}) {
	l.Error(v...)
}

func Debugf(format string, args ...interface{}) {
	l.Debugf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	l.Warnf(format, args...)
}

func Infof(format string, args ...interface{}) {
	l.Infof(format, args...)
}

func Errorf(format string, args ...interface{}) {
	l.Errorf(format, args...)
}
