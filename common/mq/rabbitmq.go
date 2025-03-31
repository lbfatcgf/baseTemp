package mq

import (
	"baseTemp/common/config"
	logger "baseTemp/common/logger"

	rbmq "github.com/rabbitmq/amqp091-go"
)

var rabbitMq *rbmq.Connection

var sendChannel *rbmq.Channel

var OnRabbitMqInit func() =func() {}
func InitRabbitMQ() {
	if config.Conf().RabbitMQ == nil {
		return
	}
	link := config.Conf().RabbitMQ.GetLink()
	conn, err := rbmq.Dial(link)
	if err != nil {
		logger.LogError("RabbitMQ connect error: %v", err)
		panic(err)
	}
	rabbitMq = conn
	sc, err := rabbitMq.Channel()
	if err != nil {
		logger.LogError("RabbitMQ create send channel error: %v", err)
		panic(err)
	}
	sendChannel = sc

	OnRabbitMqInit()
}

func CloseRabbitMQ() {
	rabbitMq.Close()
	sendChannel.Close()
}

func SetRouteQueue(exchange,routingKey, queueName string) error {
	_, err := sendChannel.QueueDeclare(
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
	err = sendChannel.QueueBind(
		queueName,                // Queue 名称
		routingKey,               // 路由键
		exchange, // Exchange 名称
		false,                    // 无需等待服务器响应
		nil,
	)
	if err != nil {

		logger.LogError("RabbitMQ bind queue error: %v", err)
		return err
	}
	return nil
}

func SetQueue(queueName string) (*rbmq.Queue, error) {

	q, err := sendChannel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		logger.LogError("RabbitMQ declare queue error: %v", err)
		return nil, err
	}

	return &q, nil
}

func SendMsg(q *rbmq.Queue, msg []byte) error {

	return sendChannel.Publish("", q.Name, false, false, rbmq.Publishing{
		ContentType: "text/plain",
		Body:        msg,
	})
}

func SendMsgToExchange(exchange,routingKey string, msg []byte) error {
	return sendChannel.Publish(
		exchange, // Exchange 名称
		routingKey,               // 路由键
		false,                    // 无需等待服务器响应
		false,                    // 无需设置消息持久化
		rbmq.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		},
	)
}
