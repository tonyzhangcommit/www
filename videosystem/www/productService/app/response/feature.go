package response

/*
	定义接口返回信息结构体
*/

import "time"

// 返回秒杀活动以及商品信息(前端展示)
type FlashEventProduct struct {
	Name           string    `json:"name"`
	Condition      string    `json:"condition"`
	ProductID      uint      `json:"productID"`
	ProductName    string    `json:"productname"`
	OriginalPrice  float64   `json:"originalprice"`
	FlashSalePrice float64   `json:"flashsaleprice"`
	Quantity       int       `json:"quantity"`
	LimitPerUser   int       `json:"limitperuser"`
	StartTime      time.Time `json:"starttime"`
	EndTime        time.Time `json:"endtime"`
	CreatedAt      time.Time `json:"createdtime"`
}

// 获取活动信息（秒杀活动请求过滤使用）
type FlashEvent struct {
	Name      string    `json:"name"`
	Condition string    `json:"condition"`
	Count     int       `json:"count"`
	StartTime time.Time `json:"starttime"`
	EndTime   time.Time `json:"endtime"`
}
