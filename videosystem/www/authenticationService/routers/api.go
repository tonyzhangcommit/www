package routers

/*
配置各个微服务api路由组
*/

import (
	"auth/global"
	"auth/request"

	"github.com/gin-gonic/gin"
)

// 用户管理服务-client路由组
func SetUserServiceClientGroupRouter(router *gin.RouterGroup) {
	clientGroup := router.Group(global.App.Config.UserServiceApi.ClientPath)
	{
		clientGroup.POST("/login", request.UserService.Login)
		clientGroup.POST("/getverifcode", request.UserService.GetVerifiCode)
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
