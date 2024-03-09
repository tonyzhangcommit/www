package management

import (
	"userservice/app/request"
	"userservice/app/response"
	"userservice/app/services"

	"github.com/gin-gonic/gin"
)

/*
	秒杀活动相关接口
*/

// 获取用户VIP类型
func GetVipType(c *gin.Context) {
	var form request.GetVipType
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if viptuype, err := services.FE.GetVipType(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, viptuype)
	}
}
