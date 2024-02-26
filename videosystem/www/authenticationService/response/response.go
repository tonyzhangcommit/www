package response

import (
	"auth/global"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
	封装相应结构体
*/

// 响应结构体
type Response struct {
	ErrorCode int         `json:"errorCode"`
	Data      interface{} `json:"data"`
	Message   string      `json:"msg"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		0,
		data,
		"success",
	})
}

func Fail(c *gin.Context, errorCode int, msg string) {
	c.JSON(http.StatusOK, Response{
		ErrorCode: errorCode,
		Data:      nil,
		Message:   msg,
	})
}

func FailExist(c *gin.Context, err global.CustomError) {
	Fail(c, err.ErrorCode, err.ErrorMsg)
}

// 内部服务错误
func LocalErrorFail(c *gin.Context, msg string) {
	Fail(c, global.Errors.LocalServiceError.ErrorCode, msg)
}

// 用户服务错误
func UserserviceFail(c *gin.Context) {
	FailExist(c,global.Errors.UserServiceError)
}

// 权限错误
func PermissionFail(c *gin.Context, msg string) {
	Fail(c, global.Errors.PermissionError.ErrorCode, msg)
}

// 非法请求
func IllegalRequestFail(c *gin.Context) {
	FailExist(c, global.Errors.IllegalRequest)
}

// 请求过于频繁
func FrequentRequestFail(c *gin.Context, msg string) {
	Fail(c, global.Errors.FrequentRequest.ErrorCode, msg)
}

// 未知错误
func UnknownErrorFail(c *gin.Context, msg string) {
	Fail(c, global.Errors.UnknownError.ErrorCode, msg)
}
