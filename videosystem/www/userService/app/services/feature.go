package services

import (
	"errors"
	"fmt"
	"userservice/app/request"
	"userservice/app/response"
	"userservice/global"
	"userservice/models"
	"userservice/utils"

	"gorm.io/gorm"
)

/*
	业务逻辑实现,为了区分不同的功能模块，这里新定义一个空结构体
*/

type feature struct {
}

var Feature = new(feature)

// 普通用户注册，如果没有指定邀请码，则默认归为默认超级管理员名下
func (Feature *feature) Register(form *request.Resister) (user models.User, err error) {
	// 验证手机号
	inputVCode := form.VarifiCode
	realVCode := utils.GetVirifCode(form.Phonenumber)
	if inputVCode != realVCode {
		err = errors.New("验证码错误")
		return
	}
	// 选定角色
	rolenameslice := global.App.Config.Roles.NameList
	rolename := ""
	if len(rolenameslice) > 3 {
		rolename = rolenameslice[2]
	} else {
		rolename = "regularUser"
	}
	var role models.Role
	if err = global.App.DB.Where("rolename = ?", rolename).First(&role).Error; err != nil {
		global.SendLogs("error", "注册失败，指定角色不存在（regularUser）", err)
		err = errors.New("指定角色不存在")

		return
	}
	// 判断用户在否存在，以手机号为验证唯一性标志
	if err = global.App.DB.Where("phonenumber = ?", form.Phonenumber).First(&models.User{}).Error; err == nil {
		err = errors.New("手机号已注册")
		return
	}
	// 判断是否存在邀请码
	manager := models.User{}
	var useractivity models.UserActivity
	if form.AgentCode != "" {
		// 判断邀请码的正确性,这里已经加了参数验证，所以这里忽略
		// agentcode 需要增加规则限制，避免SQL注入风险
		if err = global.App.DB.Where("agentcode = ?", form.AgentCode).First(&manager).Error; err != nil {
			err = errors.New("邀请码错误，注册失败")
			return
		}
		user.ParentID = &manager.ID
		user.Username = form.Name
		user.Password = utils.BcryptMake([]byte(form.Password))
		user.PhoneNumber = form.Phonenumber
		user.ParentAgentCode = manager.AgentCode
		if err = global.App.DB.Create(&user).Error; err != nil {
			global.SendLogs("error", "注册失败", err)
			err = errors.New("注册失败，请联系管理员")
		} else {
			useractivity.UserID = user.ID
			useractivity.Action = "注册"
			useractivity.Details = "新用户注册，子代理拓展"
			if erractive := global.App.DB.Create(&useractivity).Error; erractive != nil {
				global.SendLogs("error", "录入用户行为失败", erractive)
			}
		}
	} else {
		if err = global.App.DB.Where("username = ?", "desupadmin").First(&manager).Error; err != nil {
			global.SendLogs("error", "默认管理员不存在", err)
			err = errors.New("注册失败")
			return
		}
		user.ParentID = &manager.ID
		user.Username = form.Name
		user.Password = utils.BcryptMake([]byte(form.Password))
		user.PhoneNumber = form.Phonenumber
		user.ParentAgentCode = manager.AgentCode
		if err = global.App.DB.Create(&user).Error; err != nil {
			global.SendLogs("error", "注册失败", err)
			err = errors.New("注册失败，请联系管理员")

		} else {
			useractivity.UserID = user.ID
			useractivity.Action = "注册"
			useractivity.Details = "新用户注册，非代理拓展"
			if erractive := global.App.DB.Create(&useractivity).Error; erractive != nil {
				global.SendLogs("error", "录入用户行为失败", erractive)
			}
		}
	}
	return
}

// 用户名-密码登录,暂时弃用，数据库设计中允许用户名重复
func (Feature *feature) LoginByNP(form *request.LoginNP) (loginsucess response.LoginSuccess, err error) {

	return
}

// 手机号-验证码 登录
func (Feature *feature) LoginByPVC(form *request.LoginPVC) (loginsucess response.LoginSuccess, err error) {
	// 验证手机号
	inputVCode := form.VerificationCode
	realVCode := utils.GetVirifCode(form.Phonenumber)
	if inputVCode != realVCode {
		err = errors.New("验证码错误")
		return
	}
	result := global.App.DB.Where("phonenumber = ?", form.Phonenumber).First(&models.User{})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			err = errors.New("手机号未注册")
		} else {
			global.SendLogs("error", fmt.Sprintf("%s 登录失败", form.Phonenumber), err)
			err = errors.New("未知错误，登录失败")
		}
	} else {
		loginsucess.Jwt = "test jwt"
	}
	return
}
