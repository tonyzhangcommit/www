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

// 获取秒杀信息（前端展示）
func GetFEventProduct(c *gin.Context) {
	var form request.GetFlashEventProduct
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	flashepinfo, err := services.FE.GetFEventProduct(&form)
	if err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, flashepinfo)
	}
}

// 获取活动基本信息（请求过滤）
func GetFlashInfo(c *gin.Context) {
	var form request.GetFlashEvent
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	flashepinfo, err := services.FE.GetFEventInfo(&form)
	if err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, flashepinfo)
	}
}

// 更改活动商品信息（更改库存）
