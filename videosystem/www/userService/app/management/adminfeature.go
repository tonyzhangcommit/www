package management

import (
	"userservice/app/request"
	"userservice/app/response"
	"userservice/app/services"

	"github.com/gin-gonic/gin"
)

// 管理员相关
// 管理员登录,手机验证码登录
func (m adminFeature) AdminLoginByPVC(c *gin.Context) {
	var form request.LoginPVC
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	if loginS, err := services.AdminFeature.AdminLogin(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, loginS)
	}
}

// 账号密码登录
func (m adminFeature) AdminLoginByPN(c *gin.Context) {
	var form request.LoginNP
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	if loginS, err := services.AdminFeature.AdminLoginPN(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, loginS)
	}
}

// 登出
func (m adminFeature) AdminLogout(c *gin.Context) {
	var form request.LogOut
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	if err := services.AdminFeature.AdminLogout(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, "操作成功")
	}
}

// 创建管理员
func (m adminFeature) AdminCreateManager(c *gin.Context) {
	var form request.AdminCreateManager
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	} else {
		if user, err := services.AdminFeature.AdminCreateManager(&form); err != nil {
			response.BusinessFail(c, err.Error())
		} else {
			response.Success(c, user)
		}
	}
}

// 删除用户,这里只能删除下属自己的用户
func (m adminFeature) DeleteUser(c *gin.Context) {
	var form request.DeleteUser
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	} else {
		if err := services.AdminFeature.DeleteUser(&form); err != nil {
			response.BusinessFail(c, err.Error())
		} else {
			response.Success(c, nil)
		}
	}
}

// 获取当前用户列表(管理员和非管理员)
func (m adminFeature) AdminGetUserList(c *gin.Context) {
	var form request.AdminGetUserList
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	} else {
		if userlist, err := services.AdminFeature.AdminGetUserList(&form); err != nil {
			response.BusinessFail(c, err.Error())
		} else {
			response.Success(c, userlist)
		}
	}
}

// 获取单个用户信息
func (m adminFeature) AdminGetUserInfo(c *gin.Context) {
	var form request.AdminGetUserinfo
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	} else {
		if userinfo, err := services.AdminFeature.AdminGetUserInfo(&form); err != nil {
			response.BusinessFail(c, err.Error())
		} else {
			response.Success(c, userinfo)
		}
	}
}

// 编辑单个用户信息，封禁/解封
func (m adminFeature) AdminEditUserInfo(c *gin.Context) {
	var form request.AdminGetUserinfo
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	} else {
		if err := services.AdminFeature.AdminEditUserInfo(&form); err != nil {
			response.BusinessFail(c, err.Error())
		} else {
			response.Success(c, nil)
		}
	}
}

// 超管操作角色权限信息
// 获取所有角色
func (m adminFeature) GetRolesList(c *gin.Context) {
	var form request.GetRoles
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	} else {
		if roles, err := services.AdminFeature.GetRolesList(&form); err != nil {
			response.BusinessFail(c, err.Error())
		} else {
			response.Success(c, roles)
		}
	}
}

// 新增角色,传入id如果不存在则是新增，如果存在就是编辑
func (m adminFeature) AddRoles(c *gin.Context) {

}

// 删除角色
func (m adminFeature) DelRoles(c *gin.Context) {

}

// 获取所有权限
func (m adminFeature) GetPermissionList(c *gin.Context) {
	var form request.GetPermissions
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	} else {
		if roles, err := services.AdminFeature.GetPermissionList(&form); err != nil {
			response.BusinessFail(c, err.Error())
		} else {
			response.Success(c, roles)
		}
	}
}

// 新增权限
func (m adminFeature) EditPermission(c *gin.Context) {
	var form request.EditPermissions
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	} else {
		if err := services.AdminFeature.EditPermission(&form); err != nil {
			response.BusinessFail(c, err.Error())
		} else {
			response.Success(c, nil)
		}
	}
}

// 删除权限
func (m adminFeature) DelPermission(c *gin.Context) {
	var form request.GetPermissions
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	} else {
		if err := services.AdminFeature.DelPermission(&form); err != nil {
			response.BusinessFail(c, err.Error())
		} else {
			response.Success(c, nil)
		}
	}
}
