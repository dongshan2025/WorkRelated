package wsproducer

import (
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/zeromicro/go-zero/core/logx"
)

var producerPool ProducerPool

type ProducerPool struct {
	SyncProducerPool  sync.Pool
	AsyncProducerPool sync.Pool
}

func NewProducerInstance(brokers []string) {
	producerPool = ProducerPool{}

	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll       // NoResponse：不等待任何确认，直接返回；WaitForLocal：等待Leader副本确认写入后返回；WaitForAll：等待Leader副本以及相关的Follower副本确认写入后返回；
	config.Producer.Retry.Max = 5                          // 生产者重试的最大次数，默认值为3
	config.Producer.Retry.Backoff = 200 * time.Millisecond // 生产者重试之间的等待时间，默认为100毫秒
	config.Producer.Timeout = 5 * time.Second              // 超时时间

	producerPool.SyncProducerPool.New = func() any {
		config.Producer.Return.Successes = true
		producer, err := sarama.NewSyncProducer(brokers, config)
		if err != nil {
			logx.Errorf("[NewProducerInstance]Kafka同步生产者初始化失败，错误信息为：%v", err)
			panic(err)
		}

		return producer
	}

	producerPool.AsyncProducerPool.New = func() any {
		producer, err := sarama.NewAsyncProducer(brokers, config)
		if err != nil {
			logx.Errorf("[NewProducerInstance]Kafka异步生产者初始化失败，错误信息为：%v", err)
			panic(err)
		}

		return producer
	}
}

// 同步发送消息
func SyncSendMessage(topic string, message string) {
	producer := producerPool.SyncProducerPool.Get().(sarama.SyncProducer)
	defer producerPool.SyncProducerPool.Put(producer)

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		logx.Errorf("[SyncSendMessage]Kafka同步生产者发送信息失败，信息为：%s，错误信息为：%v", message, err)
		return
	}

	logx.Infof("[SyncSendMessage]Kafka同步生产者发送信息成功，信息为：%s，分区为：%d，Offset为：%d", message, partition, offset)
}

// 异步发送信息
func AsyncSendMessage(topic string, message string) {
	producer := producerPool.AsyncProducerPool.Get().(sarama.AsyncProducer)
	defer producerPool.AsyncProducerPool.Put(producer)

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	producer.Input() <- msg
	select {
	case success := <-producer.Successes():
		logx.Infof("[AsyncSendMessage]Kafka异步生产者发送信息成功，信息为：%s，分区为：%d，Offset为：%d", message, success.Partition, success.Offset)
	case err := <-producer.Errors():
		logx.Errorf("[AsyncSendMessage]Kafka异步生产者发送信息失败，信息为：%s，错误信息为：%v", message, err)
	}
}
