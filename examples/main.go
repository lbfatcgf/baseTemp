package main

import (
	"github.com/lbfatcgf/baseTemp/pkg"
	"github.com/lbfatcgf/baseTemp/pkg/config"

	// "github.com/lbfatcgf/baseTemp/pkg/db"
	"github.com/lbfatcgf/baseTemp/pkg/logger"
	// "github.com/lbfatcgf/baseTemp/pkg/mq"
)

func main() {
	pkg.ParseArgs()

	logger.InitLog(config.Conf().Mode)
	// db.Initgorm()

	// mq.InitRabbitMQ()
	// defer mq.CloseRabbitMQ()
	// 其他初始化

	pkg.AddOnStopSignal(func() {
		//添加程序退出处理
	})
	pkg.ListenerExitSignal()
}
