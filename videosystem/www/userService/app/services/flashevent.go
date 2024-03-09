package services

import (
	"errors"
	"userservice/app/request"
	"userservice/global"
	"userservice/models"

	"gorm.io/gorm"
)

type flashEvent struct{}

var FE flashEvent

func (FE flashEvent) GetVipType(form *request.GetVipType) (vtype string, err error) {
	var user models.User
	res := global.App.DB.Preload("Profile").First(&user, form.Uid)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		err = errors.New("用户不存在")
		return
	} else if res.Error != nil {
		global.SendLogs("error", "数据库查询错误", err)
		err = errors.New("未知错误")
		return
	} else {
		vtype = user.Profile.TypeVip
		return
	}

}
