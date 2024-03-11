package routers

import (
	"userservice/app/management"

	"github.com/gin-gonic/gin"
)

// 商品管理API分组
func SetProductGroupRouter(router *gin.RouterGroup) {
	router.POST("/")
}

// 商品秒杀活动API分组信息
func SetFlashPEGroupRouter(router *gin.RouterGroup) {
	router.POST("/getflasheventp", management.GetFEventProduct) // 前端展示
	router.POST("/getfeinfo", management.GetFlashInfo)          // 秒杀活动过滤
}
