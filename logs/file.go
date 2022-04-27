package logs

import (
	"io"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func DefaultFile(filename ...string) OutputOption {
	if len(filename) == 0 {
		filename = append(filename, "error.log")
	}
	hook := &lumberjack.Logger{
		Filename:   filename[0], // 日志文件路径
		MaxSize:    100,         // megabytes
		MaxBackups: 3,           // 最多保留3个备份
		MaxAge:     7,           // days
		LocalTime:  true,        // 本地时间戳
		Compress:   true,        // 是否压缩 disabled by default
	}
	return NewFileOption(zap.DebugLevel, hook)
}

func NewFileOption(level zapcore.Level, w io.Writer, config ...zapcore.EncoderConfig) OutputOption {

	if len(config) == 0 {
		cfg := zap.NewDevelopmentEncoderConfig()
		cfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		}
		config = append(config, cfg)
	}

	opt := OutputOption{
		Encoder: zapcore.NewConsoleEncoder(config[0]),
		Writer:  zapcore.AddSync(w),
		Level:   level,
	}
	return opt
}
