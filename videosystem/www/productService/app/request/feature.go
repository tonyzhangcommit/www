package request

import "time"

/*
	请求信息结构体
*/

// 获取活动&商品信息（用于前端展示）
type GetFlashEventProduct struct {
	EventId   uint `form:"eventid" json:"eventid" binding:"required"`
	ProductId uint `form:"productid" json:"productid" binding:"required"`
}

func (g GetFlashEventProduct) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"eventid.required":   "活动ID不能为空",
		"productid.required": "商品ID不能为空",
	}
}

// 获取(指定时间范围内)所有活动信息
type GetFlashEvents struct {
	StartTime time.Time `form:"starttime" json:"starttime"`
	EndTime   time.Time `form:"endtime" json:"endtime"`
}

// 获取秒杀活动信息(用于秒杀活动初步过滤)
type GetFlashEvent struct {
	EventId uint `form:"eventid" json:"eventid" binding:"required"`
}

func (g GetFlashEvent) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"eventid.required": "活动ID不能为空",
	}
}
