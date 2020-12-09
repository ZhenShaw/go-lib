package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io/ioutil"
	"os"
)

func InitLogger(logPath string, loglevel string) *zap.Logger {

	hook := lumberjack.Logger{
		Filename:   logPath, // 日志文件路径
		MaxSize:    1024,    // megabytes
		MaxBackups: 3,       // 最多保留3个备份
		MaxAge:     7,       //days
		LocalTime:  true,
		Compress:   true, // 是否压缩 disabled by default
	}
	fileWriter := zapcore.AddSync(&hook)
	consoleWriter := zapcore.Lock(os.Stdout)
	kafkaWriter := zapcore.AddSync(ioutil.Discard)

	var level zapcore.Level
	switch loglevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder //这里可以指定颜色

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	kafkaEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	cores := zapcore.NewTee(
		//打印到控制台
		zapcore.NewCore(consoleEncoder, consoleWriter, level),
		//输出到文件
		zapcore.NewCore(consoleEncoder, fileWriter, level),

		// 打印在kafka topic中（伪造的case）
		zapcore.NewCore(kafkaEncoder, kafkaWriter, level),
	)

	logger := zap.New(cores, zap.AddCaller(), zap.AddCallerSkip(0))

	//sugarLogger = logger.Sugar()

	logger.Info("DefaultLogger init success")

	return logger
}
