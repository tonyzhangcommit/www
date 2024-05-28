package routers

/*
配置各个微服务api路由组
*/

import (
	"auth/global"
	"auth/middleware"
	"auth/request"

	"github.com/gin-gonic/gin"
)

// 用户管理服务-client路由组
func SetUserServiceClientGroupRouter(router *gin.RouterGroup) {
	clientGroup := router.Group(global.App.Config.UserServiceApi.ClientPath).Use(middleware.ServiceLimit("Client", 100, 150, 200)) // 带有限流机制
	{
		platform := "app"
		clientGroup.POST("/login", func(ctx *gin.Context) { // 使用闭包进行传参
			request.UserService.Login(ctx, platform)
		})
		clientGroup.POST("/logout", middleware.JWTAUTH(platform), request.UserService.LoginOut)
		clientGroup.POST("/getverifcode", middleware.APIGetVerifCodeLimit(60), request.UserService.GetVerifiCode)
		clientGroup.POST("/register", request.UserService.Register)
		clientGroup.GET("/getuserinfo", middleware.JWTAUTH(platform), request.UserService.GetUserinfo)
		clientGroup.POST("/inproveinfo", middleware.JWTAUTH(platform), request.UserService.InproveInfo)
	}
}

// 用户管理服务-admin路由组
func SetUserServiceManageGroupRouter(router *gin.RouterGroup) {
	adminGroup := router.Group(global.App.Config.UserServiceApi.AdminPath).Use(middleware.ServiceLimit("userManager", 20, 20, 300))
	{
		platform := "management"
		adminGroup.POST("/login", func(ctx *gin.Context) {
			request.UserService.Login(ctx, platform)
		})
		adminGroup.POST("/logout", middleware.JWTAUTH(platform), request.UserService.LoginOut)
		adminGroup.POST("/getverifcode", middleware.APIGetVerifCodeLimit(120), request.UserService.GetVerifiCode)
		adminGroup.POST("/2")
		adminGroup.POST("/3")
		adminGroup.POST("/4")
		adminGroup.POST("/5")
	}
}

// 商品管理服务api管理
func SetProductServiceManageGroupRouter(router *gin.RouterGroup) {

}

// 设置订单服务api路由组
func SetOrderServiceManageGroupRouter(router *gin.RouterGroup) {}

// 秒杀活动API
func SetFlashEventServiceManageGroupRouter(router *gin.RouterGroup) {
	flashGroup := router.Group(global.App.Config.UserServiceApi.FlashEventPath)
	{
		// 这个接口需要验证管理员登录验证(预热用户信息)
		flashGroup.GET("/preuserheat", request.FlashEvent.PreheatUserInfo)
		// 测试获取用户会员信息
		flashGroup.POST("/getuserlevelinfo", request.FlashEvent.GetUserLevelInfo)
		// 活动信息预热&获取活动信息
		flashGroup.POST("/perheatproduct", request.FlashEvent.PreheatProductInfo)
		// 活动&商品展示
		flashGroup.POST("/flashepinfo", request.FlashEvent.PreheatEventProductShow)
		// 下单
		flashGroup.POST("/takeorder", middleware.JWTAUTH("app"), request.FlashEvent.PlaceOrder)
	}
}
