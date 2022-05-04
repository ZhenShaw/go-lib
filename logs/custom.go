package logs

import (
	"fmt"
	"io"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type CustomWriter struct {
	// 投递日志数据的管道，容量是batch的数倍
	logChan chan string
	// 定时写日志间隔
	interval time.Duration
	// 单次批量写日志的数量
	batch uint
	// 批量写日志缓存
	cache []string
	// 自定义的日志批量消费函数
	handle func([]string)
	// 是否阻塞写操作，默认false，当日志消费速度小于生产速度时，设置为true有可能引起日志消费不及时导致业务卡顿
	block bool
}

func NewCustomWriter(handle func([]string), batch uint, interval time.Duration) *CustomWriter {
	w := &CustomWriter{
		handle:   handle,
		interval: interval,
		logChan:  make(chan string, batch*10),
		batch:    batch,
	}
	go w.Consume()
	return w
}

func (w *CustomWriter) SetBlock(block bool) {
	w.block = block
}

func (w *CustomWriter) Consume() {
	ticker := time.NewTicker(w.interval)
	for {
		select {
		case b := <-w.logChan:
			w.cache = append(w.cache, b)
			if len(w.cache) >= int(w.batch) {
				w.doConsume()
			}
		case <-ticker.C:
			w.doConsume()
			ticker.Reset(w.interval)
		}
	}
}

func (w *CustomWriter) doConsume() {
	if len(w.cache) == 0 {
		return
	}
	w.handle(w.cache)
	w.cache = w.cache[:0]
}

func (w *CustomWriter) Write(b []byte) (int, error) {
	if !w.block && len(w.logChan) == cap(w.logChan) {
		err := fmt.Errorf("buffer overwrite, discard log: %s", string(b))
		return 0, err
	}
	w.logChan <- string(b)
	return 0, nil
}

func NewJSONOption(level zapcore.Level, w io.Writer, config ...zapcore.EncoderConfig) OutputOption {

	if len(config) == 0 {
		cfg := zap.NewProductionEncoderConfig()
		cfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		}
		config = append(config, cfg)
	}

	opt := OutputOption{
		Encoder: zapcore.NewJSONEncoder(config[0]),
		Writer:  zapcore.AddSync(w),
		Level:   level,
	}
	return opt
}
