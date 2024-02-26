package request

import (
	"auth/global"
	"auth/response"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

/*
	远程调用Userservice逻辑
*/

type uservice struct {
}

var UserService = new(uservice)

// var timeout = global.App.Config.UserServiceApi.Timeout

func (UserService *uservice) Login(c *gin.Context) {
	if c.Request.Body == nil {
		response.IllegalRequestFail(c)
		return
	}
	// 读取请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.IllegalRequestFail(c)
		return
	}
	// 创建请求客户端实例
	requester := global.NewRequestClient(3 * time.Second)
	responseBody, err := requester.DoRequest("POST", "http://127.0.0.1:19999/user/loginbypvc", map[string]string{"Content-Type": "application/json"}, body)
	if err != nil {
		// 处理转发过程中的错误
		response.UserserviceFail(c)
		return
	}

	response.Success(c, responseBody)
}
