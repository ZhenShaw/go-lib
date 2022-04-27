package logs

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewConsoleOption 默认控制台输出配置
func NewConsoleOption(level zapcore.Level, config ...zapcore.EncoderConfig) OutputOption {

	if len(config) == 0 {
		cfg := zap.NewDevelopmentEncoderConfig()
		//自定义时间格式
		cfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		}
		//输出带颜色
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config = append(config, cfg)
	}

	//使用控制台的编码格式
	consoleEncoder := zapcore.NewConsoleEncoder(config[0])

	//控制台Writer 写到标准输出，并发安全
	consoleWriter := zapcore.Lock(os.Stdout)

	opt := OutputOption{
		Encoder: consoleEncoder,
		Writer:  consoleWriter,
		Level:   level,
	}
	return opt
}
