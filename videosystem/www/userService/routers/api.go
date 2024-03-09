package routers

import (
	"net/http"
	"userservice/app/management"
	"userservice/global"

	"github.com/gin-gonic/gin"
)

// 这里主要分为两个部分：管理端/客户端
// 管理端路由组为/manage,客户端为/client
func SetClientGroupRouter(router *gin.RouterGroup) {
	router.POST("/register", management.Register)
	router.POST("/loginbypvc", management.LoginByPVC)
	router.POST("/getvarifcode", management.GetVirificationCode)
	router.GET("/getuserinfo", management.GetPersonInfo)
	router.POST("/inproveinfo", management.InprovePersonInfo)
}

// 管理端相关接口
func SetManageGroupRouter(router *gin.RouterGroup) {
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome userService Server")
		global.SendLogs("info", "这是一个测试日志。。。")
	})
}

// 秒杀活动
func SetFlashGroupRouter(router *gin.RouterGroup) {
	router.POST("/getuvip", management.GetVipType)
}
