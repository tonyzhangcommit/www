package routers

import (
	"userservice/app/management"

	"github.com/gin-gonic/gin"
)

// 这里主要分为两个部分：管理端/客户端
// 管理端路由组为/manage,客户端为/client
func SetClientGroupRouter(router *gin.RouterGroup) {
	router.POST("/register", management.UserFeature.Register)
	router.POST("/loginbypvc", management.UserFeature.LoginByPVC)
	router.POST("/getvarifcode", management.GetVirificationCode)
	router.GET("/getuserinfo", management.UserFeature.GetPersonInfo)
	router.POST("/inproveinfo", management.UserFeature.InprovePersonInfo)
}

// 管理端相关接口
// 登录，登出，创建用户，创建管理员，删除用户，封禁用户，权限管理，角色管理等..
func SetManageGroupRouter(router *gin.RouterGroup) {
	// 基础功能，登录，登出
	router.POST("/adminlogin", management.AdminFeature.AdminLoginByPN)
	router.POST("/adminlogout", management.AdminFeature.AdminLogout)
	router.POST("/createmanager", management.AdminFeature.AdminCreateManager)
	router.POST("/deletemanager", management.AdminFeature.DeleteUser)
	router.POST("/admingetuserlist", management.AdminFeature.AdminGetUserList)
	router.POST("/admingetusrinfo", management.AdminFeature.AdminGetUserInfo)
	router.POST("/adminedituserinfo", management.AdminFeature.AdminEditUserInfo)
	adminRoleGroup := router.Group("/roles")
	{
		// 增删改查,角色只能为超管
		adminRoleGroup.POST("/list", management.AdminFeature.GetRolesList)
		adminRoleGroup.POST("/editroles")
		adminRoleGroup.POST("/delrole")
	}
	adminPermissiGroup := router.Group("/permission")
	{
		// 权限的增删改查，角色只能为超管
		adminPermissiGroup.POST("/list", management.AdminFeature.GetPermissionList)
		adminPermissiGroup.POST("/editpermission",management.AdminFeature.EditPermission)
		adminPermissiGroup.POST("/delpermission",management.AdminFeature.DelPermission)
	}
}

// 秒杀活动
func SetFlashGroupRouter(router *gin.RouterGroup) {
	// 获取用户会员类型
	router.POST("/getuvip", management.GetVipType)
	// 预热用户会员信息
	router.GET("/preheat", management.PreHeat)
}
