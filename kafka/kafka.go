package kafka

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

//提供给生产者的数据
type ProduceData struct {
	Topic     string
	Partition int32
	Payload   []byte //写入队列的数据
}

//消费者读取目标并处理的参数
type ConsumeHandler struct {
	Ctx        context.Context
	Topics     []string
	GroupID    string
	ConsumeFun func(*sarama.ConsumerMessage) error //自定义的数据消费操作函数
	ErrorFun   func(*ConsumeHandler, error)        //自定义错误处理函数

	ready chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *ConsumeHandler) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as Ready
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *ConsumeHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *ConsumeHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for message := range claim.Messages() {
		session.MarkMessage(message, "")
		if err := c.ConsumeFun(message); err != nil {
			c.ErrorFun(c, err)
		} else {
			session.Commit()
		}
	}
	return nil
}

//kafka 生产者消费者实例封装
type Kafka struct {

	//kafka集群 broker连接地址
	Broker []string
	Config *sarama.Config

	//同步生产者
	producer sarama.SyncProducer
}

func DefaultKafka(broker []string) (k *Kafka, err error) {
	return NewKafka(broker, nil)
}

func NewKafka(broker []string, config *sarama.Config) (k *Kafka, err error) {

	if config == nil {
		config = sarama.NewConfig()
		// 等待服务器所有副本都保存成功后的响应
		config.Producer.RequiredAcks = sarama.WaitForAll
		// 随机的分区类型：返回一个分区器，该分区器每次选择一个随机分区
		config.Producer.Partitioner = sarama.NewRandomPartitioner
		// 是否等待成功和失败后的响应
		config.Producer.Return.Successes = true
		config.Producer.Return.Errors = true

		//提交offset的间隔时间，每秒提交一次给kafka
		config.Consumer.Offsets.CommitInterval = 1 * time.Second
	}

	if len(broker) == 0 {
		return nil, fmt.Errorf("broker is empty")
	}

	k = &Kafka{
		Broker: broker,
		Config: config,
	}

	// 使用给定代理地址和配置创建一个同步生产者
	k.producer, err = sarama.NewSyncProducer(k.Broker, k.Config)
	if err != nil {
		return
	}

	return
}

func (k *Kafka) Subscribe(h *ConsumeHandler) (err error) {
	if h == nil || h.GroupID == "" || len(h.Topics) == 0 || h.ConsumeFun == nil {
		return fmt.Errorf("illegal arguments")
	}

	h.ready = make(chan bool)

	if h.Ctx == nil {
		h.Ctx = context.Background()
	}

	if h.ErrorFun == nil {
		h.ErrorFun = func(*ConsumeHandler, error) {}
	}

	//该实例不能共享
	group, err := sarama.NewConsumerGroup(k.Broker, h.GroupID, k.Config)
	if err != nil {
		return fmt.Errorf("create consumer group failed: %w", err)
	}

	go func() {
		duration := time.Second
		for {
			// `Consume` should be called inside an infinite loop, when a server-side re-balance happens,
			// the consumer session will need to be recreated to get the new claims
			if err := group.Consume(h.Ctx, h.Topics, h); err != nil {
				err = fmt.Errorf("retry group consume in %s: %w", duration.String(), err)
				h.ErrorFun(h, err)

				//休眠重试
				time.Sleep(duration)
				duration = 2 * duration
				if duration.Minutes() > 1 {
					duration = time.Minute
				}
			}

			//重新创建，re-balance会再次触发 h.Setup()
			h.ready = make(chan bool)

			// exit: check if context was cancelled.
			if h.Ctx.Err() != nil {
				group.Close()
				return
			}
		}
	}()

	// 消费者启动超时检测
	t := time.NewTimer(k.Config.Net.DialTimeout)
	select {
	case <-h.ready:
		return nil
	case <-t.C:
		return fmt.Errorf("set up group consumer timeout")
	}
}

func (k *Kafka) Publish(data *ProduceData) (err error) {
	if data == nil {
		return fmt.Errorf("nil data")
	}

	msg := &sarama.ProducerMessage{
		Topic:     data.Topic,
		Partition: data.Partition,
		Value:     sarama.ByteEncoder(data.Payload),
	}
	_, _, err = k.producer.SendMessage(msg)

	return
}
