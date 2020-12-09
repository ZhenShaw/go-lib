package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

type OutputOption struct {
	Encoder zapcore.Encoder
	Writer  zapcore.WriteSyncer
	Level   zapcore.Level
}



func NewZapLog(caller int,outputs ...OutputOption) *zap.Logger {
	if len(outputs) == 0 {
		outputs = append(outputs, DefaultConsole())
	}
	if caller<0{
		caller=0
	}

	var cores []zapcore.Core

	for _, v := range outputs {

		core := zapcore.NewCore(v.Encoder, v.Writer, v.Level)
		cores = append(cores, core)
	}

	core := zapcore.NewTee(cores...)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(caller))


	return logger

}

//默认控制台输出配置
func DefaultConsole() OutputOption {

	//默认开发配置
	config := zap.NewDevelopmentEncoderConfig()

	//自定义时间格式
	config.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}

	//输出带颜色
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder

	//使用控制台的编码格式
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	//控制台Writer 写到标准输出，并发安全
	consoleWriter := zapcore.Lock(os.Stdout)

	opt := OutputOption{
		Encoder: consoleEncoder,
		Writer:  consoleWriter,
		Level:   zap.DebugLevel,
	}
	return opt
}

//默认文件输出配置
func DefaultFile() OutputOption {

	//默认开发配置
	config := zap.NewDevelopmentEncoderConfig()

	//自定义时间格式
	config.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}

	//输出带颜色
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder

	//使用控制台的编码格式
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	//控制台Writer 写到标准输出，并发安全
	hook := lumberjack.Logger{
		Filename:   "log.log", // 日志文件路径
		MaxSize:    100,       // megabytes
		MaxBackups: 3,         // 最多保留3个备份
		MaxAge:     7,         // days
		LocalTime:  true,      // 本地时间戳
		Compress:   true,      // 是否压缩 disabled by default
	}
	fileWriter := zapcore.AddSync(&hook)

	opt := OutputOption{
		Encoder: consoleEncoder,
		Writer:  fileWriter,
		Level:   zap.DebugLevel,
	}
	return opt
}
