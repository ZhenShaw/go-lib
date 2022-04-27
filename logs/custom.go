package logs

import (
	"io"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type CustomWriter struct {
	logChan  chan string
	interval time.Duration
	batch    uint
	cache    []string
	handle   func([]string)
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
