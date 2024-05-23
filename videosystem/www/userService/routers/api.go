package routers

import (
	"userservice/app/management"

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
// 登录，登出，创建用户，创建管理员，删除用户，封禁用户，权限管理，
func SetManageGroupRouter(router *gin.RouterGroup) {
	// 基础功能，登录，登出
	router.POST("/login")
	router.POST("/logout")
	// 用户管理路由组
	adminUserGroup := router.Group("/users")
	{
		adminUserGroup.GET("/")
		adminUserGroup.POST("/")
	}
	adminRoleGroup := router.Group("/roles")
	{
		// 增删改查
		adminRoleGroup.GET("/role")
	}
	adminPermissiGroup := router.Group("/permission")
	{
		// 增删改查,对用户增删改查
		adminPermissiGroup.POST("/")
	}
}

// 秒杀活动
func SetFlashGroupRouter(router *gin.RouterGroup) {
	// 获取用户会员类型
	router.POST("/getuvip", management.GetVipType)
	// 预热用户会员信息
	router.GET("/preheat", management.PreHeat)
}
