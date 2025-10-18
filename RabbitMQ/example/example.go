package example

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

// 消息发布
func PublishMessage(channel *amqp.Channel, exchange, routingKey, body string) error {
	return channel.Publish(
		exchange,   // 指定交换机，若为""表示默认交换机
		routingKey, // 路由键，根据交换机类型决定消息怎么路由
		false,      // mandatory，是否强制投递，若为true且无法路由到队列，则会触发Basic.Return（需要监听返回）
		false,      // immediate，是否立即投递，很少使用，RabbitMQ通常不支持，建议设为false
		amqp.Publishing{ // 消息及其元数据
			ContentType:  "text/plain",    // 内容类型 纯文本:text/plain JSON数据:application/json HTML文档:text/html JPEG图片:image/jpeg
			Body:         []byte(body),    // 消息内容
			DeliveryMode: amqp.Persistent, // 持久化消息
		},
	)
}

// mandatory=true时，表示“找不到接收方要告诉我”（确保消息不被悄悄丢掉）
// immediate=true时，表示“没有消费者就别投了”（对方不在线就别发），RabbitMQ已经不支持该参数为true了，应该一直设置为false

// 带mandatory回退处理机制的Rabbit生产者示例
func MandatoryProducerExample() {
	// 连接到RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer channel.Close()

	// 声明交换机（topic类型）
	err = channel.ExchangeDeclare(
		"my-exchange",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	// 设置return回退监听（必须在Publish之前设置）
	returns := channel.NotifyReturn(make(chan amqp.Return))

	// 模拟发送消息，但没有任何队列绑定这个key，消息将被回退
	err = channel.Publish(
		"my-exchange",
		"unmatched-key",
		true,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte("This message will be returned"),
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	// 检查是否被回退（注意：这是异步的）
	select {
	case ret := <-returns:
		fmt.Printf("消息被回退，原因：%s，交换机：%s，路由键：%s，内容：%s", ret.ReplyText, ret.Exchange, ret.RoutingKey, string(ret.Body))
	case <-time.After(2 * time.Second):
		fmt.Println("消息已成功路由（没有被回退）")
	}
}

// 创建AMQP连接
func CreateConnection(amqpURI string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// 监听连接关闭事件
	go func() {
		err := <-conn.NotifyClose(make(chan *amqp.Error))
		log.Printf("连接关闭：%v", err)
	}()

	return conn, nil
}

// 创建和使用通道
func UseChannel(conn *amqp.Connection) error {
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	// 启用确认模式
	if err := channel.Confirm(false); err != nil {
		return err
	}

	confirms := channel.NotifyPublish(make(chan amqp.Confirmation))
	go func(ch chan amqp.Confirmation) {
		select {
		case con := <-ch:
			fmt.Printf("Ack: %t, DeliveryTag: %d", con.Ack, con.DeliveryTag)
		default:
			fmt.Println("default")
		}
	}(confirms)

	return nil
}

// ================================================== 可靠的消息生产者 ==================================================
type ReliableProducer struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	confirms chan amqp.Confirmation
}

func NewReliableProducer(amqpURI string) (*ReliableProducer, error) {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// 启用发布者确认
	if err = channel.Confirm(false); err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}

	confirms := channel.NotifyPublish(make(chan amqp.Confirmation))

	return &ReliableProducer{
		conn:     conn,
		channel:  channel,
		confirms: confirms,
	}, nil
}

func (p *ReliableProducer) Publish(exchange, routingKey string, message []byte) error {
	return p.channel.Publish(
		exchange,
		routingKey,
		false,

		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         message,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
}

// ================================================== 高性能消息消费者 ==================================================
type BatchConsumer struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	deliverCh <-chan amqp.Delivery
	batchSize int
	timeout   time.Duration
}

func NewBatchConsumer(amqpURI string, queueName string, batchSize int, timeout time.Duration) (*BatchConsumer, error) {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// 设置Qos预取计数
	if err := channel.Qos(batchSize, 0, false); err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}

	deliveries, err := channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}

	return &BatchConsumer{
		conn:      conn,
		channel:   channel,
		deliverCh: deliveries,
		batchSize: batchSize,
		timeout:   timeout,
	}, nil
}

func (c *BatchConsumer) ConsumeBatch(handler func([]amqp.Delivery) error) {
	var batch []amqp.Delivery
	timer := time.NewTicker(c.timeout)

	for {
		select {
		case delivery, ok := <-c.deliverCh:
			if !ok {
				// 通道关闭，处理剩余批次
				if len(batch) > 0 {
					handler(batch)
				}
				return
			}

			batch = append(batch, delivery)
			if len(batch) >= c.batchSize {
				if err := handler(batch); err != nil {
					// 处理错误，可能需要重试或记录
					fmt.Printf("批处理错误：%v", err)
				}
				batch = nil
				timer.Reset(c.timeout)
			}
		case <-timer.C:
			if len(batch) > 0 {
				if err := handler(batch); err != nil {
					fmt.Printf("超时批处理时发生错误：%v", err)
				}
				batch = nil
			}
			timer.Reset(c.timeout)
		}
	}
}

// ================================================== 高级特性与最佳实践 ==================================================
// 确认模式			描述						   适用场景
// 自动确认			消息一点送达立即确认			测试环境，可接受消息丢失
// 手动确认			显示调用Ack/Nack			   生成环境，要求可靠传递
// 批量确认			批量处理后再确认			    高性能场景，减少网络开销

// 手动确认示例
func HandleDelivery(delivery amqp.Delivery) {
	defer func() {
		if err := recover(); err != nil {
			// 处理异常，拒绝消息并重新入队
			delivery.Nack(false, true)
			log.Printf("消息处理异常：%v", err)
		}
	}()

	// 业务逻辑处理
	if processMessage(delivery.Body) {
		delivery.Ack(false) // 单条确认
	} else {
		delivery.Nack(false, true) // 拒绝并重新入队
	}
}

func processMessage(message []byte) bool {
	fmt.Println(string(message))
	return true
}

// ========== 连接恢复策略 ==========
type ConnectionManager struct {
	amqpURI     string
	conn        *amqp.Connection
	reconnectCh chan struct{}
	maxRetries  int
	retryDelay  time.Duration
}

func NewConnectionManager(amqpURI string, maxRetries int, retryDelay time.Duration) *ConnectionManager {
	mgr := &ConnectionManager{
		amqpURI:     amqpURI,
		reconnectCh: make(chan struct{}, 1),
		maxRetries:  maxRetries,
		retryDelay:  retryDelay,
	}

	go mgr.monitorConnection()

	return mgr
}

func (m *ConnectionManager) monitorConnection() {
	for range m.reconnectCh {
		m.reconnectWithRetry()
	}
}

func (m *ConnectionManager) reconnectWithRetry() {
	for i := 0; i < m.maxRetries; i++ {
		conn, err := amqp.Dial(m.amqpURI)
		if err == nil {
			m.conn = conn
			log.Println("连接恢复成功")
			return
		}

		log.Printf("连接尝试 %d 失败：%v", i+1, err)
		time.Sleep(m.retryDelay * time.Duration(i+1))
	}

	log.Println("达到最大重试次数，连接恢复失败")
}

// ========== 连接池管理 ==========
type ConnectionPool struct {
	amqpURI string
	pool    chan *amqp.Connection
	maxSize int
	mu      sync.Mutex
}

func NewConnectionPool(amqpURI string, maxSize int) *ConnectionPool {
	pool := &ConnectionPool{
		amqpURI: amqpURI,
		pool:    make(chan *amqp.Connection, maxSize),
		maxSize: maxSize,
	}

	// 预热连接池
	for i := 0; i < maxSize/2; i++ {
		conn, err := amqp.Dial(amqpURI)
		if err == nil {
			pool.pool <- conn
		}
	}

	return pool
}

func (p *ConnectionPool) Get() (*amqp.Connection, error) {
	select {
	case conn := <-p.pool:
		return conn, nil
	default:
		return amqp.Dial(p.amqpURI)
	}
}

func (p *ConnectionPool) Put(conn *amqp.Connection) {
	select {
	case p.pool <- conn:
	default:
		conn.Close()
	}
}

// ========== 批量发布优化 ==========
type Message struct {
	RoutingKey string
	Content    string
}

func BatchPublish(channel *amqp.Channel, exchange string, messages []Message) error {
	// 启用事务模式
	if err := channel.Tx(); err != nil {
		return err
	}

	for _, msg := range messages {
		if err := channel.Publish(
			exchange,
			msg.RoutingKey,
			false,
			false,
			amqp.Publishing{
				Body:         []byte(msg.Content),
				DeliveryMode: amqp.Persistent,
			},
		); err != nil {
			channel.TxRollback()
			return err
		}
	}

	// 提交事务
	return channel.TxCommit()
}
