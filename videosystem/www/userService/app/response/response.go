package response

import (
	"net/http"
	"userservice/global"

	"github.com/gin-gonic/gin"
)

/*
	封装response
*/

// 相应结构体
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

// FailByError 失败响应 返回自定义错误的错误码、错误信息
func FailByError(c *gin.Context, error global.CustomError) {
	Fail(c, error.ErrorCode, error.ErrorMsg)
}

// ValidateFail 请求参数验证失败
func ValidateFail(c *gin.Context, msg string) {
	Fail(c, global.Errors.ValidateError.ErrorCode, msg)
}

// BusinessFail 业务逻辑失败
func BusinessFail(c *gin.Context, msg string) {
	Fail(c, global.Errors.BusinessError.ErrorCode, msg)
}
