package request

/*
	定义部分请求的请求体
*/

// 秒杀活动下单
type TakeFlashOrder struct {
	EventID   uint `json:"eventid" binding:"required"`
	ProductID uint `json:"peoductid" binding:"required"`
	UserID    uint `json:"userid" binding:"required"`
	Count     int  `json:"count" binding:"required"`
}

func (t TakeFlashOrder) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"eventid.required":   "活动ID不能为空",
		"peoductid.required": "产品ID不能为空",
		"userid.required":    "用户ID不能为空",
		"count.required":     "商品数量不能为空",
	}
}
