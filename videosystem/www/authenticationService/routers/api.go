package routers

/*
配置各个微服务api路由组
*/

import (
	"auth/global"
	"auth/middleware"
	"auth/request"
	"auth/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 用户管理服务-client路由组
func SetUserServiceClientGroupRouter(router *gin.RouterGroup) {
	router.GET("/test", func(ctx *gin.Context) {
		secretKey, err := utils.GenerateSecretKey(32)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"err": "error",
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"key": secretKey,
			})
		}
	})
	clientGroup := router.Group(global.App.Config.UserServiceApi.ClientPath).Use(middleware.ServiceLimit("Client", 100, 150, 200))
	{
		clientGroup.POST("/login", request.UserService.Login)
		clientGroup.POST("/logout", request.UserService.Login)
		clientGroup.POST("/getverifcode", middleware.APIGetVerifCodeLimit(6), request.UserService.GetVerifiCode)
		clientGroup.POST("/register", request.UserService.Register)
		clientGroup.GET("/getuserinfo", request.UserService.GetUserinfo)
		clientGroup.POST("/inproveinfo", request.UserService.InproveInfo)
	}
}

// 用户管理服务-admin路由组
func SetUserServiceManageGroupRouter(router *gin.RouterGroup) {
	adminGroup := router.Group(global.App.Config.UserServiceApi.AdminPath).Use(middleware.ServiceLimit("userManager", 20, 20, 300))
	_ = adminGroup
}

// 商品管理服务api管理
func SetProductServiceManageGroupRouter(router *gin.RouterGroup) {}

// 设置订单服务api路由组
func SetOrderServiceManageGroupRouter(router *gin.RouterGroup) {}

// 秒杀活动API
func SetFlashEventServiceManageGroupRouter(router *gin.RouterGroup) {
	flashGroup := router.Group(global.App.Config.UserServiceApi.FlashEventPath)
	{
		// 这个接口需要验证管理员登录验证
		flashGroup.GET("/preuserheat", request.FlashEvent.PreheatUserInfo)
		flashGroup.POST("/getuserlevelinfo", request.FlashEvent.GetUserLevelInfo)
		// 商品信息预热
		flashGroup.GET("/preproductheat", request.FlashEvent.PreheatProductInfo)
		// 下单
		flashGroup.POST("/placeorder", request.FlashEvent.PlaceOrder)
	}

}
