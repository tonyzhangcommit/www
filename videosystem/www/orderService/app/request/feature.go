package request

/*
	请求结构体
*/

// 秒杀活动下单
type TakeFlashOrder struct {
	EventID   uint `json:"eventid" binding:"require"`
	ProductID uint `json:"peoductid" binding:"require"`
	UserID    uint `json:"userid" binding:"require"`
	Count     int  `json:"count" binding:"require"`
}

func (t TakeFlashOrder) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"eventid.required":   "活动ID不能为空",
		"peoductid.required": "产品ID不能为空",
		"userid.required":    "用户ID不能为空",
		"count.required":     "商品数量不能为空",
	}
}

// 普通下单
type TakeRegularOrder struct {
	ProductsID []uint `json:"productsid" binding:"require"`
	UserID     uint   `json:"userid" binding:"require"`
	Count      int    `json:"count" binding:"require"`
}

// 修改订单 。。。
