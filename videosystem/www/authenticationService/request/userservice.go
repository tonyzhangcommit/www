package request

import (
	"auth/global"
	"auth/middleware"
	"auth/response"
	"auth/utils"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
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
	body, err := GetRequestBody(c)
	if err != nil {
		response.IllegalRequestFail(c)
		return
	}
	// 创建请求客户端实例
	requester := global.NewRequestClient(time.Duration(global.App.Config.UserServiceApi.Timeout) * time.Second)
	// 创建请求
	req, err := http.NewRequest("POST", remoteurl, bytes.NewBuffer(body))
	if err != nil {
		go global.SendLogs("error", "创建请求错误", err)
		return
	}
	// 设置请求头
	for key, value := range map[string]string{"Content-Type": "application/json"} {
		req.Header.Set(key, value)
	}
	// 发出请求
	resp, err := requester.Client.Do(req)
	if err != nil {
		go global.SendLogs("error", "发送请求错误", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		go global.SendLogs("error", "解析response错误", err)
		return
	}

	// 序列化
	var loginres response.LoginRes
	var commonres response.Response
	err = json.Unmarshal(responseBody, &commonres)
	if err != nil {
		response.UserserviceFail(c)
		return
	}
	// 判断是否登录成功，登录成功生成jwt认证
	if commonres.ErrorCode == 0 {
		err = json.Unmarshal(responseBody, &loginres)
		if err != nil {
			response.UserserviceFail(c)
			return
		}
		jwtout, err := response.JwtService.CreateJwtToken("app", strconv.Itoa(int(loginres.Data.ID)), loginres.Data.Roles)
		if err != nil {
			go global.SendLogs("error", "生成令牌出错", err)
			response.LocalErrorFail(c, "生成令牌出错")
			return
		}
		response.Success(c, jwtout)
	} else {
		response.Success(c, commonres)
	}
}

// 退出登录
func (UserService *uservice) LoginOut(c *gin.Context) {
	if err := middleware.JoinBlackList(c.Keys["token"].(*jwt.Token)); err != nil {
		response.JwtTokenErrorFail(c, "登出错误，请联系管理员")
	}
	response.Success(c, nil)
}

// 获取验证码
func (UserService *uservice) GetVerifiCode(c *gin.Context) {
	remoteurl := utils.JoinStrings(global.App.Config.UserServiceApi.BaseUrl, global.App.Config.UserServiceApi.ClientUrl.Getverifcode)
	PostRequest(c, global.App.Config.UserServiceApi.Timeout, remoteurl)
}

// 注册
func (UserService *uservice) Register(c *gin.Context) {
	remoteurl := utils.JoinStrings(global.App.Config.UserServiceApi.BaseUrl, global.App.Config.UserServiceApi.ClientUrl.Register)
	PostRequest(c, global.App.Config.UserServiceApi.Timeout, remoteurl)
}

// 获取个人信息
func (UserService *uservice) GetUserinfo(c *gin.Context) {
	remoteurl := utils.JoinStrings(global.App.Config.UserServiceApi.BaseUrl, global.App.Config.UserServiceApi.ClientUrl.GetuserInfo)
	GetRequest(c, global.App.Config.UserServiceApi.Timeout, remoteurl)
}

// 注册
func (UserService *uservice) InproveInfo(c *gin.Context) {
	remoteurl := utils.JoinStrings(global.App.Config.UserServiceApi.BaseUrl, global.App.Config.UserServiceApi.ClientUrl.InproveInfo)
	PostRequest(c, global.App.Config.UserServiceApi.Timeout, remoteurl)
}
