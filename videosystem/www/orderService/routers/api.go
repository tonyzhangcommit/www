package routers

import (
	"userservice/app/management"

	"github.com/gin-gonic/gin"
)

func SetOrderGroupRouter(router *gin.RouterGroup) {
	// 秒杀活动下单
	router.POST("/takeflashorder", management.TakeFlashOrder)
	// websocket 连接
	router.GET("/ws",)
}
