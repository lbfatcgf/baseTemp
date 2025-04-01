package mq

import (
	"baseTemp/common/config"
	logger "baseTemp/common/logger"
	"context"
	"fmt"

	rbmq "github.com/rabbitmq/amqp091-go"
)

var rabbitMq = make(map[string]*rbmq.Connection)

var sendChannel = make(map[string]*rbmq.Channel)
var channelNames = make([]string, 0)
var OnRabbitMqInit func() = func() {}

var currentChannel int = 0
var consumeMsgChan = make(map[string]chan rbmq.Delivery)

var consumeCtx = make(map[string]context.CancelFunc)

func InitRabbitMQ() {
	if config.Conf().RabbitMQ == nil {
		return
	}
	parentLink := config.Conf().RabbitMQ.GetLink()
	conn, err := rbmq.Dial(parentLink)
	if err != nil {
		logger.LogError("RabbitMQ connect error: %v", err)
		panic(err)
	}
	sc, err := conn.Channel()
	if err != nil {
		logger.LogError("RabbitMQ create send channel error: %v", err)
		panic(err)
	}
	channelNames = append(channelNames, config.Conf().RabbitMQ.Host+config.Conf().RabbitMQ.User)
	rabbitMq[config.Conf().RabbitMQ.Host+config.Conf().RabbitMQ.User] = conn
	sendChannel[config.Conf().RabbitMQ.Host+config.Conf().RabbitMQ.User] = sc
	logger.LogInfo(fmt.Sprintf("RabbitMQ %s connect success", config.Conf().RabbitMQ.Host))
	connNode()

	OnRabbitMqInit()
}

func connNode() {
	if config.Conf().RabbitMQ.Node != nil && len(*config.Conf().RabbitMQ.Node) > 0 {
		for _, node := range *config.Conf().RabbitMQ.Node {
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
			logger.LogInfo(fmt.Sprintf("RabbitMQ %s connect success", node.Host))
			channelNames = append(channelNames, node.Host+node.User)
			rabbitMq[node.Host+node.User] = nodeConn
			sendChannel[node.Host+node.User] = nodeSc
		}
	}
}

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

func SetRouteQueue(exchange, routingKey, queueName string) error {
	var channel *rbmq.Channel
	for _, sc := range sendChannel {
		if sc.IsClosed() {
			continue
		}
		channel = sc
		break
	}
	_, err := channel.QueueDeclare(
		queueName, // Queue 名称
		true,      // 持久化
		false,     // 自动删除
		false,     // 不等待服务器响应
		false,     // 无额外参数
		nil,
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
		channelNames = append(channelNames[:currentChannel], channelNames[currentChannel+1:]...)
		return NextChannel()
	}
	next := sendChannel[channelNames[currentChannel]]
	fmt.Println(channelNames[currentChannel])
	currentChannel++
	return next, nil
}
func SetQueue(queueName string) (*rbmq.Queue, error) {
	channel, err := NextChannel()
	if err != nil {
		return nil, err
	}

	q, err := channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		logger.LogError("RabbitMQ declare queue error: %v", err)
		return nil, err
	}

	return &q, nil
}

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

func StopConsumeMsg(queue string) {
	if cancel, ok := consumeCtx[queue]; ok {
		cancel()
	}
	delete(consumeCtx, queue)
}

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
				continue
			}
			for {
				if qch[queue] == nil {
					msg = nil
					break
				}
				m, ok := <-msg
				if !ok {
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
