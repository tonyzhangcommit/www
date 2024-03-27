package management

import (
	"fmt"
	"userservice/app/request"
	"userservice/app/response"
	"userservice/app/services"

	"github.com/gin-gonic/gin"
)

/*
	视图函数
*/

// 秒杀活动下单
func TakeFlashOrder(c *gin.Context) {
	var form request.TakeFlashOrder
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	res, err := services.TakeFlashOrder(&form)
	if err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, res)
	}
}

func Test(c *gin.Context) {
	fmt.Println(services.WebSocketclients)
	go services.WebsocketSendMessage("3", "处理您的请求时出现问题，请检查订单详情或联系客服。")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// 普通下单
func TakeOrder(c *gin.Context) {

}

// 修改订单
func EditOrder(c *gin.Context) {

}

// 取消订单
func CancelOrder(c *gin.Context) {

}

// 删除订单
func DeleteOrder(c *gin.Context) {

}

// 查看订单
func GetOrder(c *gin.Context) {

}

// 申请退款
func RefundOrder(c *gin.Context) {

}
