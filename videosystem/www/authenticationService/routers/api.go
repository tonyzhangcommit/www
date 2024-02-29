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
		//
		clientGroup.POST("/login", request.UserService.Login)
		clientGroup.POST("/getverifcode", middleware.APIGetVerifCodeLimit(6), request.UserService.GetVerifiCode)
		clientGroup.POST("/register", request.UserService.Register)
	}
}

// 用户管理服务-admin路由组
func SetUserServiceManageGroupRouter(router *gin.RouterGroup) {
	adminGroup := router.Group(global.App.Config.UserServiceApi.AdminPath)
	_ = adminGroup
}

// 用户管理服务-admin路由组

// 商品管理服务api管理
