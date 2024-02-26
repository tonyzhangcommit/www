package request

import (
	"auth/global"
	"auth/utils"

	"github.com/gin-gonic/gin"
)

/*
	远程调用Userservice逻辑
*/

type uservice struct {
}

var UserService = new(uservice)

// 登录
func (UserService *uservice) Login(c *gin.Context) {
	remoteurl := utils.JoinStrings(global.App.Config.UserServiceApi.BaseUrl, global.App.Config.UserServiceApi.ClientUrl.Login)
	BaseRequest(c, global.App.Config.UserServiceApi.Timeout, remoteurl)
}

// 获取验证码
func (UserService *uservice) GetVerifiCode(c *gin.Context) {
	remoteurl := utils.JoinStrings(global.App.Config.UserServiceApi.BaseUrl, global.App.Config.UserServiceApi.ClientUrl.Getverifcode)
	BaseRequest(c, global.App.Config.UserServiceApi.Timeout, remoteurl)
}

// 注册
func (UserService *uservice) Register(c *gin.Context) {
	remoteurl := utils.JoinStrings(global.App.Config.UserServiceApi.BaseUrl, global.App.Config.UserServiceApi.ClientUrl.Register)
	BaseRequest(c, global.App.Config.UserServiceApi.Timeout, remoteurl)
}

