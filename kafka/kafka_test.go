package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/zhenshaw/go-lib/logs"
	"testing"
	"time"
)

func TestDefaultKafka(t *testing.T) {

	//开启生产者和消费者实例
	k, err := DefaultKafka([]string{"localhost:9092"})
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	h := &ConsumeHandler{
		GroupID:    "mykafka",
		Topics:     []string{"mykafka"},
		ConsumeFun: customMsg,
		ErrorFun:   errHandler,
	}

	err = k.Subscribe(h)
	if err != nil {
		logs.Error(err.Error())
		return
	}

	for i := 0; i < 5; i++ {
		data := &ProduceData{
			Topic:   "mykafka",
			Payload: []byte("get one message: " + fmt.Sprint(i)),
		}

		err = k.Publish(data)
		if err != nil {
			logs.Error(err.Error())
			return
		}

		time.Sleep(1 * time.Second)
	}

	time.Sleep(10 * time.Second)

}

func customMsg(msg *sarama.ConsumerMessage) (err error) {

	if msg == nil {
		return
	}

	logs.Info("消费了 %s:%d : %s", msg.Topic, msg.Offset, string(msg.Value))

	return nil

}

func errHandler(h *ConsumeHandler, err error) {

	if err == nil {
		return
	}
	logs.Error("%s发生错误：%s", h.Topics, err.Error())
}
