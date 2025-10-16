package pkg

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/lbfatcgf/baseTemp/pkg/logger"
	"github.com/lbfatcgf/baseTemp/pkg/mq"
)

var stopCallbackList = make([]func(), 0)

// AddOnStopSignal 添加停止信号回调函数
func AddOnStopSignal(callback func()) {
	stopCallbackList = append(stopCallbackList, callback)
}

// StopSingalHandler 停止信号处理函数
func StopSingalHandler() {
	for _, callback := range stopCallbackList {
		callback()
	}
}

//阻塞监听退出信号
func ListenerExitSignal() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		StopSingalHandler()
		mq.CloseRabbitMQ()
		logger.CloseLog()
		os.Exit(0)
	}()
}
