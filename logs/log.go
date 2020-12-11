package logs

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

var (
	z *zap.Logger

	filepath = "log.log"
	caller   = 1

	openConsole   = true
	openFile      = true
	consoleOutput = DefaultConsole()
	fileOutput    = DefaultFile(filepath)
	outputOptions []OutputOption
)

func GetLogger() *zap.Logger {
	return z
}

func AddOutput(output OutputOption) {
	outputOptions = append(outputOptions, output)
	initLog()
}

func SetCaller(c int) {
	caller = c
	initLog()
}

func SetLogPath(path string) {
	path = strings.TrimSuffix(path, "/")
	filepath = path
	fileOutput = DefaultFile(filepath)
	initLog()
}

func CloseConsoleOutput() {
	openConsole = false
	initLog()
}

func CloseFileOutput() {
	openFile = false
	initLog()
}

func initLog() {
	outputs := outputOptions
	if openConsole {
		outputs = append(outputs, consoleOutput)
	}
	if openFile {
		outputs = append(outputs, fileOutput)
	}
	z = NewZapLog(caller, outputs...)
	//l = z.Sugar()
}

func init() {
	initLog()
}

func Debug(f interface{}, v ...interface{}) {
	z.Debug(formatLog(f, v...))
}

func Warn(f interface{}, v ...interface{}) {
	z.Warn(formatLog(f, v...))
}

func Info(f interface{}, v ...interface{}) {
	z.Info(formatLog(f, v...))
}

func Error(f interface{}, v ...interface{}) {
	z.Error(formatLog(f, v...))
}

func formatLog(f interface{}, v ...interface{}) string {
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
