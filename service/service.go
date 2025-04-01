package service

import (
	"baseTemp/common/logger"
	"baseTemp/common/mq"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	ramp "github.com/rabbitmq/amqp091-go"
)

var q *ramp.Queue
var count int = 0

func InitChiTongDaService(engine *gin.Engine) {
	cq, err := mq.SetQueue("quorum_test_q")
	if err != nil {
		logger.LogError("queue init error", err)
	}
	q = cq
	engine.Match([]string{"GET", "POST"}, "/baseTemp", func(ctx *gin.Context) {
		if q != nil {
			mq.SendMsg(q, []byte(fmt.Sprintf("dd:%d", count)))
		}
		count++
		ctx.JSON(200, gin.H{
			"message": "Hello, world!",
		})
	})
	engine.GET("stop",func(ctx *gin.Context) {
		mq.StopConsumeMsg("quorum_test_q")
		ctx.JSON(200, gin.H{
			"message": "stop consume",
		})
	})
	listen()
}

func listen() {
	ctx:= context.Background()
	msgs := mq.ConsumeMsg(
		ctx,
		"quorum_test_q",
		"test_consumer",
		true,
		false,
		false,
		false,
		nil)
	go func(ch <-chan ramp.Delivery) {
		for {
			msg, ok := <-ch
			if ok {
				fmt.Println(string(msg.Body))
			}
		}
	}(msgs)
}
