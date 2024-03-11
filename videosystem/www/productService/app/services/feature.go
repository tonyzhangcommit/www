package services

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"userservice/app/request"
	"userservice/app/response"
	"userservice/global"
	"userservice/models"

	"gorm.io/gorm"
)

/*
	业务逻辑实现,为了区分不同的功能模块，这里新定义一个空结构体
*/

type flashEvent struct{}

var FE = new(flashEvent)

// 获取活动以及商品信息
func (f *flashEvent) GetFEventProduct(form *request.GetFlashEventProduct) (fpinfo response.FlashEventProduct, err error) {
	var result struct {
		models.FlashSaleEvent
		models.FlashSaleEventProduct
	}
	// 查询指定数据
	res := global.App.DB.Table("flasheventproduct").Select("flasheventproduct.*, flashsaleevent.*").Joins("join flashsaleevent on flashsaleevent.id = flasheventproduct.eventid").Where("flasheventproduct.eventid = ? AND flasheventproduct.productid = ?", form.EventId, form.ProductId).First(&result)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		err = errors.New("活动不存在")
		return
	} else if res.Error != nil {
		global.SendLogs("error", "mysql 查询活动信息报错", err)
		err = errors.New("内部错误")
		return
	}
	// 查询商品名称
	var product = models.Membership{}
	if err = global.App.DB.First(&product, result.FlashSaleEventProduct.ProductID).Error; err != nil {
		global.SendLogs("error", "mysql 查询商品名称报错", err)
		err = errors.New("内部错误")
		return
	}
	fpinfo.Name = result.FlashSaleEvent.Name
	fpinfo.Condition = result.FlashSaleEvent.Condition
	fpinfo.ProductID = result.FlashSaleEventProduct.ProductID
	fpinfo.ProductName = product.Name
	fpinfo.OriginalPrice = result.FlashSaleEventProduct.OriginalPrice
	fpinfo.FlashSalePrice = result.FlashSaleEventProduct.FlashSalePrice
	fpinfo.Quantity = result.FlashSaleEventProduct.Quantity
	fpinfo.LimitPerUser = result.FlashSaleEventProduct.LimitPerUser
	fpinfo.StartTime = result.FlashSaleEvent.StartTime
	fpinfo.EndTime = result.FlashSaleEvent.EndTime
	fpinfo.CreatedAt = result.FlashSaleEvent.CreatedAt
	return
}

// 获取活动基本信息（使用redis缓存查询结果）
func (f *flashEvent) GetFEventInfo(form *request.GetFlashEvent) (finfo response.FlashEvent, err error) {
	key := "flash:base:info:" + strconv.Itoa(int(form.EventId))
	if val, errredis := global.App.Redis.Get(context.Background(), key).Result(); errredis != nil {
		var event = models.FlashSaleEvent{}
		if err = global.App.DB.First(&event, form.EventId).Error; err != nil {
			global.SendLogs("error", "mysql 查询活动报错", err)
			err = errors.New("内部错误")
			return
		}
		finfo.Name = event.Name
		finfo.Condition = event.Condition
		finfo.StartTime = event.StartTime
		finfo.EndTime = event.EndTime
		if eventJSON, err := json.Marshal(event); err != nil {
			global.SendLogs("error", "序列化活动基本信息报错", err)
		} else {
			if errsredis := global.App.Redis.Set(context.Background(), key, eventJSON, 0).Err(); errsredis != nil {
				global.SendLogs("error", "redis 设置商品活动基本信息报错", errredis)
			}
		}
	} else {
		err = json.Unmarshal([]byte(val), &finfo)
		if err != nil {
			global.SendLogs("error", "redis 活动基本信息序列化到结构体报错", errredis)
			err = errors.New("内部错误")
		}
	}
	return
}
