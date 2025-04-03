package main

import (
	"baseTemp/cmd"
	"baseTemp/common/config"
	"baseTemp/common/db"
	logger "baseTemp/common/logger"
	"baseTemp/common/mq"
	"baseTemp/service"
	"baseTemp/tools"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)


func main() {
	time.LoadLocation("Asia/Shanghai")
	if cmd.ParseArgs() {
		return
	}

	if config.Conf().Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	logger.InitLog(config.Conf().Mode)
	mq.OnRabbitMqInit = func() {

	}
	mq.InitRabbitMQ()
	defer mq.CloseRabbitMQ()
	db.Initgorm()
	db.MigrateTable()
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有源
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	service.InitChiTongDaService(router)

	tools.AddOnStopSignal(func() {
		fmt.Println("stop")
	})
	listenerSignal()
	router.Run(":" + *cmd.Port)
	
}

func listenerSignal() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		tools.StopSingalHandler()
		mq.CloseRabbitMQ()
		os.Exit(0)
	}()
}
func MigrateTable() ( ) {
	if !config.Conf().InitTable {
		return
	}
	err := db.DB().AutoMigrate()
	if err != nil {
		panic(err.Error())
	}
}
