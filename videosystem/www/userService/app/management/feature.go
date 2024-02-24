package management

import (
	"time"
	"userservice/app/request"
	"userservice/app/response"
	"userservice/app/services"
	"userservice/global"
	"userservice/utils"

	"github.com/gin-gonic/gin"
)

/*
	视图函数
*/

// 注册,用户名，手机号，密码，必填项
func Register(c *gin.Context) {
	var form request.Resister
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if user, err := services.Feature.Register(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, user)
	}
}

// 登录 用户名-密码， 暂时弃用
func LoginByNP(c *gin.Context) {
	var form request.LoginNP
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if loginS, err := services.Feature.LoginByNP(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, loginS)
	}
}

// 登录 手机号-验证码
func LoginByPVC(c *gin.Context) {
	var form request.LoginPVC
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if loginS, err := services.Feature.LoginByPVC(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, loginS)
	}
}

// 编辑个人信息
func InprovePersonInfo() {

}

// 退出登录
func Logout() {

}

// 发送验证码
func GetVirificationCode(c *gin.Context) {
	var form request.GetVirifCode
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
	} else {
		varifCode := utils.GenerateNumberCode(6)
		// 首先存放至内存,以5分钟为有效期
		utils.SetVirifCode(form.Phonenumber, varifCode, 5*time.Minute)
		// 这里调用短信接口发送短信，暂时省略
		go global.SendLogs("info", form.Phonenumber+"验证码："+varifCode)
		response.Success(c, "验证码已发送")
	}
}
