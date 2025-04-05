package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/lbfatcgf/baseTemp/common/logger"
	"github.com/lbfatcgf/baseTemp/common/mq"
	"github.com/lbfatcgf/baseTemp/tools"
)

func main() {
	

}

func listenerSignal() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		tools.StopSingalHandler()
		mq.CloseRabbitMQ()
		logger.CloseLog()
		os.Exit(0)
	}()
}
