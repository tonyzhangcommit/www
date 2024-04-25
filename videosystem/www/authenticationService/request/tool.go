package request

import (
	"auth/global"
	"auth/response"
	"auth/utils"
	"errors"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取请求体，仅适用POST请求
func GetRequestBody(c *gin.Context) (body []byte, err error) {
	if c.Request.Body == nil {
		err = errors.New("非法请求")
		return
	}
	// 读取请求体
	body, err = io.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	return
}

// 封装POST请求函数
func PostRequest(c *gin.Context, timeout int, remoteurl string) {
	body, err := GetRequestBody(c)
	if err != nil {
		response.IllegalRequestFail(c)
		return
	}
	// 创建请求客户端实例
	requester := global.NewRequestClient(time.Duration(timeout) * time.Second)
	responseBody, err := requester.DoRequest("POST", remoteurl, map[string]string{"Content-Type": "application/json"}, body)
	if err != nil {
		// 处理转发过程中的错误
		response.UserserviceFail(c)
		go global.SendLogs("error", "转发错误", err)
		return
	}
	response.Success(c, responseBody)
}

// 封装GET请求
func GetRequest(c *gin.Context, timeout int, remoteurl string) {
	queryParams := c.Request.URL.Query().Encode()
	remoteurl = utils.JoinStrings(remoteurl, "?", queryParams)
	// 创建请求客户端实例
	requester := global.NewRequestClient(time.Duration(timeout) * time.Second)
	responseBody, err := requester.DoRequest("GET", remoteurl, map[string]string{}, []byte{})
	if err != nil {
		// 处理转发过程中的错误
		response.UserserviceFail(c)
		go global.SendLogs("error", "转发错误", err)
		return
	}
	response.Success(c, responseBody)
}
