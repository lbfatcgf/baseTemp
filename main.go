package main

import (
	"baseTemp/cmd"
	"baseTemp/common/config"
	"baseTemp/common/db"
	logger "baseTemp/common/logger"
	"baseTemp/common/mq"
	"baseTemp/service"
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	time.LoadLocation("Asia/Shanghai")
	if(cmd.ParseArgs()){
		return
	}
	config.InitConfig(*cmd.ConfigPath)
	
	if config.Conf().Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	logger.InitLog(config.Conf().Mode)
	mq.OnRabbitMqInit=func() {

	}
	mq.InitRabbitMQ()
	defer mq.CloseRabbitMQ()
	db.Initgorm()
	db.MigrateTable()
	router := gin.Default()
	// 创建会话存储
	// store := cookie.NewStore([]byte("secret")) // 这里的 "secret" 应该是一个安全的密钥
	// router.Use(sessions.Sessions("session", store))
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有源
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	service.InitChiTongDaService(router)
	router.Run(":" + *cmd.Port)
}
