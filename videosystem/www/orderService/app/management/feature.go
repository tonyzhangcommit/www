package management

import (
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
