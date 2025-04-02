package service

import (
	"github.com/gin-gonic/gin"
)



func InitChiTongDaService(engine *gin.Engine) {
	
	engine.Match([]string{"GET", "POST"}, "/baseTemp", func(ctx *gin.Context) {
	
		ctx.JSON(200, gin.H{
			"message": "Hello, world!",
		})
	})
	
}
