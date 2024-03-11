package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
	"userservice/app/request"
	"userservice/global"
	"userservice/models"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type flashEvent struct{}

var FE = new(flashEvent)

func (f *flashEvent) GetVipType(form *request.GetVipType) (vtype string, err error) {
	var user models.User
	// 首先查询redis
	key := "flash:UserID:" + strconv.FormatUint(uint64(form.Uid), 10)
	val, err := global.App.Redis.Get(context.Background(), key).Result()
	if err == redis.Nil || err != nil {
		global.SendLogs("info", "redis 获取用户会员类型失败")
		res := global.App.DB.Preload("Profile").First(&user, form.Uid)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			err = errors.New("用户不存在")
			return
		} else if res.Error != nil {
			global.SendLogs("error", "数据库查询错误", err)
			err = errors.New("未知错误")
			return
		} else {
			if time.Now().Before(user.Profile.ExpVipDate) {
				vtype = user.Profile.TypeVip
			} else {
				vtype = "regularUser"
			}
			return
		}
	} else {
		vtype = val
		return
	}
}

// 活动信息预热
func (f *flashEvent) PreHeat() (err error) {
	type userid struct {
		ID      uint   `gorm:"column:id"`
		TypeVip string `gorm:"column:typevip"`
	}
	var users []userid
	now := time.Now()
	res := global.App.DB.Debug().Table("users").Select("users.id", "profiles.typevip").Joins("left join profiles on profiles.user_id = users.id").Where("profiles.expvipdate > ? AND users.isbanned = ?", now, false).Scan(&users)
	if res.Error != nil {
		global.SendLogs("error", "mysql获取有效会员信息失败", err)
		err = errors.New("未知错误")
	} else {
		fmt.Println("------------->")
		fmt.Println(users)
		fmt.Println("------------->")
		for _, item := range users {
			key := "flash:UserID:" + strconv.FormatUint(uint64(item.ID), 10)
			value := item.TypeVip
			dur := time.Hour * 24
			err = global.App.Redis.Set(context.Background(), key, value, dur).Err()
			if err != nil {
				global.SendLogs("error", "redis保存用户信息失败", err)
				err = errors.New("未知错误")
			}
		}
	}
	return
}
