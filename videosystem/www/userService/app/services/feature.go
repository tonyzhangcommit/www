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

type feature struct{}

var Feature = new(feature)

// 普通用户注册，如果没有指定邀请码，则默认归为默认超级管理员名下
func (Feature *feature) Register(form *request.Resister) (user models.User, err error) {
	// 验证手机号
	inputVCode := form.VarifiCode
	realVCode := utils.GetVirifCode(form.Phonenumber)
	// 暂时去除验证码验证功能
	_ = inputVCode
	_ = realVCode
	// if inputVCode != realVCode {
	// 	err = errors.New("验证码错误")
	// 	return
	// }
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
	if err = global.App.DB.Where("username = ?", form.Name).First(&models.User{}).Error; err == nil {
		err = errors.New("用户名已注册")
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
		user.Roles = []models.Role{role}
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
		user.Roles = []models.Role{role}
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
func (Feature *feature) LoginByNP(form *request.LoginNP) (user response.LoginRes, err error) {

	return
}

// 手机号-验证码 登录
func (Feature *feature) LoginByPVC(form *request.LoginPVC) (loginRes response.LoginRes, err error) {
	// 验证手机号
	inputVCode := form.VerificationCode
	realVCode := utils.GetVirifCode(form.Phonenumber)
	if inputVCode != realVCode {
		err = errors.New("验证码错误")
		return
	}
	var user models.User
	err = global.App.DB.Preload("Roles").Where("phonenumber = ?", form.Phonenumber).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("手机号未注册")
		} else {
			global.SendLogs("error", fmt.Sprintf("%s 登录失败", form.Phonenumber), err)
			err = errors.New("未知错误，登录失败")
		}
		return
	}
	var roles []string
	for _, user := range user.Roles {
		roles = append(roles, user.RoleName)
	}
	loginRes.AgentCode = user.AgentCode
	loginRes.CreatedAt = user.CreatedAt
	loginRes.UpdatedAt = user.UpdatedAt
	loginRes.Username = user.Username
	loginRes.Roles = roles
	loginRes.ParentID = *user.ParentID
	loginRes.ID = user.ID
	loginRes.Isbanned = user.IsBanned
	loginRes.Phonenumber = user.PhoneNumber
	return
}

// 获取个人信息，请求参数为手机号或者用户名（目的是防止get请求被猜测）
func (Feature *feature) GetPersonInfo(form *request.GetPersonInfo) (userinfo response.UserInfo, err error) {
	if form.Name == "" && form.Phonenumber == "" {
		err = errors.New("缺少参数")
		return
	}
	var user = models.User{}
	if form.Name != "" {
		err = global.App.DB.Preload("Roles").Preload("Profile").Where("username=?", form.Name).First(&user).Error
	} else {
		err = global.App.DB.Preload("Roles").Preload("Profile").Where("phonenumber=?", form.Phonenumber).First(&user).Error
	}
	if err != nil {
		err = errors.New("用户不存在")
		return
	}

	var roles = []string{}
	for _, item := range user.Roles {
		roles = append(roles, item.RoleName)
	}
	userinfo.Username = user.Username
	userinfo.PhoneNumber = user.PhoneNumber
	userinfo.Address = user.Profile.Address
	userinfo.Identification = user.Profile.Identification
	userinfo.Email = user.Profile.Email
	userinfo.VIP = user.Profile.VIP
	userinfo.TypeVip = user.Profile.TypeVip
	userinfo.ExpVipDate = user.Profile.ExpVipDate
	userinfo.Preferences = user.Profile.Preferences
	userinfo.Sex = user.Profile.Sex
	userinfo.Roles = roles
	userinfo.AgentCode = user.AgentCode
	userinfo.ParentAgentCode = user.ParentAgentCode
	userinfo.IsBanned = user.IsBanned
	return
}

// 完善个人信息
func (Feature *feature) InprovePersonInfo(form *request.InproveInfo) (err error) {
	fmt.Println(form)
	if err = global.App.DB.First(&models.User{}, form.UserID).Error; err != nil {
		err = errors.New("用户不存在")
		return
	}
	var userprofile = models.Profile{}
	if err = global.App.DB.Where("user_id=?", form.UserID).First(&userprofile).Error; err != nil {
		global.SendLogs("error", "查询用户信息报错：mysql 报错：", err)
		err = errors.New("内部错误！")
		return
	}
	userprofile.Address = form.Address
	userprofile.Sex = form.Sex
	userprofile.Identification = form.Identification
	userprofile.Email = form.Email
	userprofile.Preferences = form.Preferences
	if err = global.App.DB.Save(&userprofile).Error; err != nil {
		global.SendLogs("error", "更新用户信息报错：mysql 报错：", err)
		err = errors.New("内部错误！")
		return
	}
	return
}

// 退出登录
func (Feature *feature) Logout(form *request.GetPersonInfo) (err error) {
	return
}
