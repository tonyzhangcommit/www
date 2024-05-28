package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"
	"userservice/app/request"
	"userservice/app/response"
	"userservice/global"
	"userservice/models"
	"userservice/utils"

	"gorm.io/gorm"
)

const (
	SuAdminRoleName  = "superAdmin"
	AdminRoleName    = "admin"
	RegisterName     = "注册"
	LoginName        = "登录"
	LogoutName       = "登出"
	CreateManager    = "创建管理员"
	IllegalLoginName = "非法登录"
)

// 管理员登录,手机验证码
func (m adminFeature) AdminLogin(form *request.LoginPVC) (loginRes response.LoginRes, err error) {
	formCode := form.VerificationCode
	realCode := utils.GetVirifCode(form.Phonenumber)
	if formCode != realCode {
		err = errors.New("验证码错误，登录失败")
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
	// 判断是否封禁
	if user.IsBanned {
		err = errors.New("用户已封禁")
		return
	}

	// 判断是否为管理员
	if !isAdmin(user) {
		// 这种非法登录的可能性大，需要检查接口
		err = errors.New("非法登录")
		var useractivity models.UserActivity
		useractivity.Action = IllegalLoginName
		useractivity.UserID = user.ID
		if err := global.App.DB.Create(&useractivity).Error; err != nil {
			global.SendLogs("error", "入库用户行为失败", err)
		}
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
	loginRes.ID = user.ID
	loginRes.Isbanned = user.IsBanned
	loginRes.Phonenumber = user.PhoneNumber
	return
}

// 账号密码登录
func (m adminFeature) AdminLoginPN(form *request.LoginNP) (loginRes response.LoginRes, err error) {
	var user models.User
	err = global.App.DB.Preload("Roles").Where("username = ?", form.Name).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("用户不存在")
		} else {
			global.SendLogs("error", fmt.Sprintf("%s 登录失败", form.Name), err)
			err = errors.New("未知错误，登录失败")
		}
		return
	}

	// 判断状态，是否封禁
	if user.IsBanned {
		err = errors.New("用户已封禁，登陆失败")
		return
	}

	// 验证密码
	if !utils.BcryptMakeCheck([]byte(form.Password), user.Password) {
		err = errors.New("密码错误")
		return
	}

	// 判断是否为管理员
	if !isAdmin(user) {
		// 这种非法登录的可能性大，需要检查接口
		err = errors.New("非法登录")
		var useractivity models.UserActivity
		useractivity.Action = IllegalLoginName
		useractivity.UserID = user.ID
		if err := global.App.DB.Create(&useractivity).Error; err != nil {
			global.SendLogs("error", "入库用户行为失败", err)
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
	loginRes.ID = user.ID
	loginRes.Isbanned = user.IsBanned
	loginRes.Phonenumber = user.PhoneNumber
	return
}

// 判断用户是否为管理员
func isAdmin(user models.User) bool {
	for _, item := range user.Roles {
		if item.RoleName == SuAdminRoleName || item.RoleName == AdminRoleName {
			return true
		}
	}
	return false
}

// 判断用户是否为超管
func isSuAdmin(user models.User) bool {
	for _, item := range user.Roles {
		if item.RoleName == SuAdminRoleName {
			return true
		}
	}
	return false
}

// 登出
func (m adminFeature) AdminLogout(form *request.LogOut) (err error) {
	var user models.User

	if err = global.App.DB.Where("uid=", form.UserId).First(&user).Error; err != nil {
		err = errors.New("非法请求，用户不存在")
		global.SendLogs("error", "非法退出登录")
		return
	}
	var useractivity models.UserActivity
	useractivity.Action = LogoutName
	useractivity.UserID = user.ID
	if err := global.App.DB.Create(&useractivity).Error; err != nil {
		global.SendLogs("error", "入库用户行为失败", err)
	}
	return
}

// create manager
func (m adminFeature) AdminCreateManager(form *request.AdminCreateManager) (user models.User, err error) {
	var curuser models.User
	if err = global.App.DB.Preload("Roles").First(&curuser, form.Uid).Error; err != nil {
		err = errors.New("用户不存在，非法请求")
		return
	}
	// 判断是否为管理员
	if !isAdmin(curuser) {
		// 这种非法登录的可能性大，需要检查接口
		err = errors.New("非法操作")
		global.SendLogs("error", "非法创建管理员")
		return
	}
	// 获取role 对象
	var managerrole models.Role
	if err = global.App.DB.Where("rolename=?", AdminRoleName).First(&managerrole).Error; err != nil {
		err = errors.New("内部错误，角色不存在")
		global.SendLogs("error", "创建manager失败，原因是获取角色对象失败", err)
		return
	}
	// 开始创建管理员
	user.ParentID = &curuser.ID
	user.Username = form.UserName
	user.Password = utils.BcryptMake([]byte(form.Password))
	user.PhoneNumber = form.Phonenumber
	user.Roles = []models.Role{managerrole}
	user.ParentAgentCode = curuser.AgentCode
	user.AgentCode = utils.GenerateRCode(8)
	if err = global.App.DB.Create(&user).Error; err != nil {
		return
	} else {
		// 增加日志,
		var useractivity models.UserActivity
		useractivity.Action = CreateManager
		useractivity.UserID = curuser.ID
		if err := global.App.DB.Create(&useractivity).Error; err != nil {
			global.SendLogs("error", fmt.Sprintf("超管%d创建管理员%d同步日志记录失败", curuser.ID, user.ID), err)
		}
	}
	return
}

// delete user or manager
func (m adminFeature) DeleteUser(form *request.DeleteUser) (err error) {
	var adminuser models.User
	var sonuser models.User
	// 首先判断用户是否存在
	if err = global.App.DB.First(&adminuser, form.Uid).Error; err != nil {
		err = errors.New("uid 不存在")
		return
	}
	if err = global.App.DB.First(&sonuser, form.DeleteUid).Error; err != nil {
		err = errors.New("deleteuid 不存在")
		return
	}

	// 判断两者的层级
	if !isSonUser(adminuser, sonuser) {
		err = errors.New("权限错误，非法请求")
		return
	}
	// 开始删除，第一步接触关联
	global.App.DB.Model(&sonuser).Association("Roles").Clear()
	// 删除用户
	err = global.App.DB.Delete(&sonuser).Error

	return
}

// 获取当前用户列表
func (m adminFeature) AdminGetUserList(form *request.AdminGetUserList) (res response.UsersPagesInfo, err error) {

	var usersinfo []models.User

	if form.PageSize <= 0 {
		form.PageSize = 10
	}

	if form.CurrentPage <= 0 {
		form.CurrentPage = 1
	}

	// 判断角色
	var operateUser models.User
	if err = global.App.DB.Preload("Roles").First(&operateUser, form.Uid).Error; err != nil {
		err = errors.New("用户不存在")
		return
	}
	if !isAdmin(operateUser) {
		// 这种非法登录的可能性大，需要检查接口
		err = errors.New("非法请求")
		return
	}
	// 逻辑过程，首先将获取到的用户信息列表保存到缓存中，并设置一定的过期时间，然后根据不同的请求参数，返回不同的用户信息
	// redis 键值  users:manager:getuserlist:pagesize:10
	usersCacheKey := "users:manager:getuserlist:pagesize:" + strconv.FormatInt(form.PageSize, 10)
	cacheuserinfo, err := global.App.Redis.Get(context.Background(), usersCacheKey).Result()
	if err != nil || form.IsFirstRequest {
		// 删除缓存
		global.App.Redis.Del(context.Background(), usersCacheKey)
		// 表示取缓存失败或者为第一次（或者用户点击刷新）时，此时需要读取数据库
		if isSuAdmin(operateUser) {
			err = global.App.DB.Preload("Roles", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "rolename", "createat")
			}).Preload("Profile").Select("id", "username", "phonenumber", "parentid", "agentcode", "parentagentcode", "isbanned", "createat").Find(&usersinfo).Error
			if err != nil {
				return
			}
		} else {
			usersinfo, err = getSubordinates(uint64(form.Uid))
			if err != nil {
				err = errors.New("查询错误，请联系管理员")
				return
			}
		}
		// 设置缓存，将userlist设置为缓存
		jsonData, masherr := json.Marshal(usersinfo)
		if masherr != nil {
			global.SendLogs("error", "序列化usersinfo失败：", err)
			err = errors.New("内部错误，请联系管理员")
			return
		}

		if err = global.App.Redis.Set(context.Background(), usersCacheKey, jsonData, time.Hour*2).Err(); err != nil {
			global.SendLogs("error", "设置usersinfo缓存失败：", err)
			err = errors.New("内部错误，请联系管理员")
			return
		}

	} else {
		err = json.Unmarshal([]byte(cacheuserinfo), &usersinfo)
		if err != nil {
			global.SendLogs("error", "反序列化usersinfo失败：", err)
			err = errors.New("内部错误，请联系管理员")
			return
		}
	}
	// 开始分页
	res = cutPageDeal(form.PageSize, form.CurrentPage, usersinfo)
	return
}

// 获取单个用户信息
func (m adminFeature) AdminGetUserInfo(form *request.AdminGetUserinfo) (user models.User, err error) {
	var muser models.User
	if err = global.App.DB.First(&muser, form.Uid).Error; err != nil {
		err = errors.New("uid 不存在")
		return
	}
	if err = global.App.DB.Select("id", "username", "phonenumber", "parentid", "agentcode", "parentagentcode", "isbanned", "createat").Preload("Roles").Preload("Profile").First(&user, form.TargetUid).Error; err != nil {
		err = errors.New("用户不存在")
		return
	}
	if !isSonUser(muser, user) {
		global.SendLogs("error", fmt.Sprintf("用户 %d 非法查询目标用户 %d 信息", muser.ID, user.ID))
		err = errors.New("非法操作")
	}
	return
}

// 编辑单个用户信息,这里为是封禁用户
func (m adminFeature) AdminEditUserInfo(form *request.AdminGetUserinfo) (err error) {
	var muser, user models.User
	if err = global.App.DB.First(&muser, form.Uid).Error; err != nil {
		err = errors.New("uid 不存在")
		return
	}
	if err = global.App.DB.First(&user, form.TargetUid).Error; err != nil {
		err = errors.New("用户不存在")
		return
	}
	if !isSonUser(muser, user) {
		global.SendLogs("error", fmt.Sprintf("用户 %d 非法操作目标用户 %d 信息", muser.ID, user.ID))
		err = errors.New("非法操作")
		return
	}
	// 开始操作  0封禁，1表示解禁
	user.IsBanned = !form.Action
	if err = global.App.DB.Save(&user).Error; err != nil {
		global.SendLogs("error", fmt.Sprintf("用户 %d 解/封禁目标用户 %d 信息失败", muser.ID, user.ID), err)
		err = errors.New("操作失败")
	}

	return
}

// 两个参数，判断第二个用户是否为第一个用户的子用户、
func isSonUser(father models.User, son models.User) bool {
	// 超级管理员的情况
	if father.ParentAgentCode == "" {
		return true
	}
	// 开始判断
	rawfatheragentcode := father.AgentCode
	if son.ParentAgentCode == rawfatheragentcode {
		return true
	} else {
		sonAgentCode := son.ParentAgentCode
		if sonAgentCode == "" {
			return false
		}
		for father.AgentCode != son.ParentAgentCode {
			// 往上翻一级
			if err := global.App.DB.Where("agentcode=?", sonAgentCode).First(&father).Error; err != nil {
				return false
			}
			if rawfatheragentcode == father.ParentAgentCode {
				return true
			}
			sonAgentCode = father.ParentAgentCode
		}
	}
	return false
}

// 获取管理员下所有的用户
func getSubordinates(uid uint64) ([]models.User, error) {
	var result []models.User
	var children []models.User

	err := global.App.DB.Where("parentid = ?", uid).Find(&children).Error
	if err != nil {
		return nil, err
	}

	result = append(result, children...)

	for _, child := range children {
		subChildren, err := getSubordinates(uint64(child.ID))
		if err != nil {
			return nil, err
		}
		result = append(result, subChildren...)
	}

	return result, nil
}

// 分页处理函数，
func cutPageDeal(pagesize int64, currentpage int64, users []models.User) (upi response.UsersPagesInfo) {

	totalUsers := int64(len(users)) // 总人数
	if totalUsers == 0 {
		return
	}
	totalpage := int64(math.Ceil(float64(totalUsers) / float64(pagesize))) // 总页数
	if currentpage > totalpage {
		currentpage = totalpage
	}

	start := (currentpage - 1) * pagesize
	end := start + pagesize
	if end > totalUsers {
		end = totalUsers
	}
	upi.CurrentPage = currentpage
	upi.Pagesize = pagesize
	upi.Total = totalUsers
	upi.Users = users[start:end]
	return
}

// 超管操作角色权限信息
