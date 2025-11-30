package wsconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"ws/internal/global"
	"ws/internal/svc"

	"github.com/IBM/sarama"
	"github.com/zeromicro/go-zero/core/logx"
)

// 主要用于转发通知时使用
var svcCtx *svc.ServiceContext

func NewConsumer(svc *svc.ServiceContext, brokers []string) {
	svcCtx = svc
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Return.Errors = true                  // 消费时将错误信息写入到Errors通道
	config.Consumer.Fetch.Default = 3 * 1024 * 1024       // 默认请求的字节数
	config.Consumer.Offsets.Initial = sarama.OffsetNewest // 从最新的offset读取，如果设置为OffsetOldest则从最旧的offset读取
	config.Consumer.Offsets.AutoCommit.Enable = false     // 设置为手动提交，默认为自动提交true

	// 连接Kafka服务器，可以传入多个broker，用逗号连接
	consumer, err := sarama.NewConsumerGroup(brokers, "NotifyConsumerGroup", config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	// 消费消息
	for {
		ctx := context.Background()
		err := consumer.Consume(ctx, []string{global.WS_NOTIFYTOPIC}, &NotifyConsumerGroupHandler{})
		if err != nil {
			fmt.Printf("consume error: %v", err)
		}
	}
}

type NotifyConsumerGroupHandler struct {
}

func (c *NotifyConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *NotifyConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *NotifyConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		notify := &global.NotifyMessage{}
		err := json.Unmarshal(msg.Value, notify)
		if err != nil {
			logx.Errorf("[ConsumeClaim]反序列化时发生错误，错误信息为：%v", err)
			return err
		}

		for _, wsid := range notify.Wsid {
			cli := svcCtx.Manager.GetClient(wsid)
			if cli != nil {
				cli.Egress <- []byte(notify.Message)
				logx.Infof("[ConsumeClaim]消费者处理了[%s]发送的消息[%s]", wsid, notify.Message)
			} else {
				logx.Errorf("[ConsumeClaim]消费者没有发现[%s]的客户端，消息[%s]发送失败", wsid, notify.Message)
			}
		}

		logx.Infof("[ConsumeClaim]消费者处理[%v]的消息[%s]完成", notify.Wsid, notify.Message)

		// 设置标记
		sess.MarkMessage(msg, "")
		// 手动提交
		sess.Commit()
	}

	return nil
}
