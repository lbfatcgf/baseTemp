package mq

import (
	"context"
	"fmt"
	"slices"

	"github.com/lbfatcgf/baseTemp/pkg/config"
	"github.com/lbfatcgf/baseTemp/pkg/logger"
	rbmq "github.com/rabbitmq/amqp091-go"
)

var rabbitMq = make(map[string]*rbmq.Connection)

var sendChannel = make(map[string]*rbmq.Channel)
var channelNames = make([]string, 0)

// OnRabbitMqInit 回调函数，在初始化rabbitMq连接后时调用
var OnRabbitMqInit func() = func() {}

// currentChannel 当前使用的channel
var currentChannel int = 0

// consumeMsgChan 消费消息通道
var consumeMsgChan = make(map[string]chan rbmq.Delivery)

// consumeCtx 消费上下文
var consumeCtx = make(map[string]context.CancelFunc)

// InitRabbitMQ 初始化rabbitMq连接
func InitRabbitMQ() {
	if config.Conf().RabbitMQ == nil || len(*config.Conf().RabbitMQ) == 0 {
		return
	}

	connNode()

	OnRabbitMqInit()
}

// connNode 获取rabbitMq连接
func connNode() {
	connCount := 0
	for _, node := range *config.Conf().RabbitMQ {
		nodeLink := node.GetLink()
		nodeConn, err := rbmq.Dial(nodeLink)
		if err != nil {
			logger.LogError("RabbitMQ connect error: %v", err)
			continue
		}
		nodeSc, err := nodeConn.Channel()
		if err != nil {
			logger.LogError("RabbitMQ create send channel error: %v", err)
			continue
		}
		connCount++
		logger.LogInfo(fmt.Sprintf("RabbitMQ %s connect success", node.Host))
		channelNames = append(channelNames, node.Host+node.User)
		rabbitMq[node.Host+node.User] = nodeConn
		sendChannel[node.Host+node.User] = nodeSc
	}
	if connCount == 0 {
		panic("RabbitMQ connect error")
	}
}

// CloseRabbitMQ 关闭rabbitMq连接
func CloseRabbitMQ() {

	for _, conn := range rabbitMq {
		conn.Close()
	}
	for _, sc := range sendChannel {
		sc.Close()
	}
	for k, ctx := range consumeCtx {
		ctx()
		delete(consumeCtx, k)
	}

	for k, ch := range consumeMsgChan {
		close(ch)
		delete(consumeMsgChan, k)
	}
}

// SetRouteQueueDefault 设置默认队列绑定关系
func SetRouteQueueDefault(exchange, routingKey, queueName string) error {
	return SetRouteQueue(exchange, routingKey, queueName, true, false, false, false, nil)
}

// SetRouteQueue 设置队列绑定关系
func SetRouteQueue(exchange, routingKey, queueName string, durable bool, autoDelete bool, exclusive bool, noWait bool, args rbmq.Table) error {
	var channel *rbmq.Channel
	for _, sc := range sendChannel {
		if sc.IsClosed() {
			continue
		}
		channel = sc
		break
	}
	_, err := channel.QueueDeclare(
		queueName,  // Queue 名称
		durable,    // 是否持久化
		autoDelete, // 是否自动删除
		exclusive,  // 是否排他
		noWait,     // 是否不等待服务器响应
		args,       // 其他参数
	)
	if err != nil {

		logger.LogError("RabbitMQ bind queue error: %v", err)
		return err
	}
	err = channel.QueueBind(
		queueName,  // Queue 名称
		routingKey, // 路由键
		exchange,   // Exchange 名称
		false,      // 无需等待服务器响应
		nil,
	)
	if err != nil {

		logger.LogError("RabbitMQ bind queue error: %v", err)
		return err
	}
	return nil
}

// NextChannel 获取下一个channel
func NextChannel() (*rbmq.Channel, error) {

	if len(channelNames) != len(rabbitMq) || len(channelNames) != len(sendChannel) {

	}
	if len(channelNames) == 0 || len(rabbitMq) == 0 || len(sendChannel) == 0 {
		return nil, fmt.Errorf("RabbitMQ channel is not exist")
	}
	if currentChannel >= len(channelNames) && currentChannel > 0 {
		currentChannel = len(channelNames) % currentChannel
	}

	if sendChannel[channelNames[currentChannel]] == nil || sendChannel[channelNames[currentChannel]].IsClosed() {

		if rabbitMq[channelNames[currentChannel]] != nil {
			rabbitMq[channelNames[currentChannel]].Close()
		} else {
			delete(rabbitMq, channelNames[currentChannel])
		}
		if sendChannel[channelNames[currentChannel]] != nil {
			sendChannel[channelNames[currentChannel]].Close()
		} else {
			delete(sendChannel, channelNames[currentChannel])
		}

		// channelNames = append(channelNames[:currentChannel], channelNames[currentChannel+1:]...)
		if len(channelNames) > 1 {

			channelNames = slices.Delete(channelNames, currentChannel, currentChannel+1)
		} else {
			channelNames = channelNames[:0]
		}
		return NextChannel()
	}
	next := sendChannel[channelNames[currentChannel]]
	fmt.Println(channelNames[currentChannel])
	currentChannel++
	return next, nil
}

// SetQueueDefault 设置默认队列

func SetQueueDefault(queueName string) (*rbmq.Queue, error) {
	return SetQueue(queueName, true, false, false, false, nil)
}

// SetQueue 设置队列
func SetQueue(queueName string, durable bool, autoDelete bool, exclusive bool, noWait bool, args rbmq.Table) (*rbmq.Queue, error) {
	channel, err := NextChannel()
	if err != nil {
		return nil, err
	}

	q, err := channel.QueueDeclare(queueName, true, false, false, false, args)
	if err != nil {
		logger.LogError("RabbitMQ declare queue error: %v", err)
		return nil, err
	}

	return &q, nil
}

// SendMsg 发送消息
func SendMsg(q *rbmq.Queue, msg []byte) error {
	channel, err := NextChannel()
	if err != nil {
		return err
	}
	return channel.Publish("", q.Name, false, false, rbmq.Publishing{
		ContentType: "text/plain",
		Body:        msg,
	})
}

// SendMsgToExchange 发送消息到exchange
func SendMsgToExchange(exchange, routingKey string, msg []byte) error {
	channel, err := NextChannel()
	if err != nil {
		return err
	}
	return channel.Publish(
		exchange,   // Exchange 名称
		routingKey, // 路由键
		false,      // 无需等待服务器响应
		false,      // 无需设置消息持久化
		rbmq.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		},
	)
}

// StopConsumeMsg 停止消费消息
func StopConsumeMsg(queue string) {
	if cancel, ok := consumeCtx[queue]; ok {
		cancel()
	}
	delete(consumeCtx, queue)
}

// ConsumeMsg 开始消费消息
func ConsumeMsg(ctx context.Context, queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args rbmq.Table) <-chan rbmq.Delivery {
	allChan := make(chan rbmq.Delivery, 1000)
	consumeMsgChan[queue] = allChan
	nctx, cancel := context.WithCancel(ctx)
	consumeCtx[queue] = cancel
	go func(ctx context.Context, qch map[string]chan rbmq.Delivery) {
		for {
			if ctx.Err() != nil {
				break
			}
			if qch[queue] == nil {
				break
			}
			ch, err := NextChannel()
			if err != nil {
				break
			}
			msg, err := ch.ConsumeWithContext(ctx, queue, consumer, autoAck, exclusive, noLocal, noWait, args)
			if err != nil {
				logger.LogDebug("RabbitMQ consume error: %v", err)
				continue
			}
			if config.Conf().Mode != "release" {
				logger.LogDebug("RabbitMQ consume msg: %v", channelNames[currentChannel])
			}
			for {
				if qch[queue] == nil {
					msg = nil
					break
				}
				m, ok := <-msg
				if !ok {
					if config.Conf().Mode != "release" {
						logger.LogDebug("RabbitMQ consume channel closed")
					}
					break
				}
				allChan <- m
			}
		}
		if qch[queue] != nil {
			close(qch[queue])
			delete(qch, queue)
		}

	}(nctx, consumeMsgChan)

	return allChan
}
