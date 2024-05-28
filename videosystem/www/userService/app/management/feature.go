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
	视图函数，这里主要用于参数验证
*/

// 根据业务对方法进行区分
type userFeature struct{}
type adminFeature struct{}

// 外部调用
var UserFeature = new(userFeature)
var AdminFeature = new(adminFeature)

// 注册,用户名，手机号，密码，必填项
func (u userFeature) Register(c *gin.Context) {
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
func (u userFeature) LoginByNP(c *gin.Context) {
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
func (u userFeature) LoginByPVC(c *gin.Context) {
	var form request.LoginPVC
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if user, err := services.Feature.LoginByPVC(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, user)
	}
}

// 获取个人信息
func (u userFeature) GetPersonInfo(c *gin.Context) {
	var form request.GetPersonInfo
	if err := c.ShouldBindQuery(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if userinfo, err := services.Feature.GetPersonInfo(&form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, userinfo)
	}
}

// 编辑（完善）个人信息
func (u userFeature) InprovePersonInfo(c *gin.Context) {
	var form request.InproveInfo
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	if err := services.Feature.InprovePersonInfo(&form); err != nil {
		response.BusinessFail(c, "更新失败")
	} else {
		response.Success(c, "编辑/完善资料成功")
	}
}

// 更新（忘记）密码
func (u userFeature) UpdatePwd(c *gin.Context) {

}

// 开通会员,分为开通会员和升级会员
func (u userFeature) OpenVip(c *gin.Context) {

}

// 发送验证码（通用方法）
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

