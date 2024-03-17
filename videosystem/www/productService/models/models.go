/*
这里是数据库表设计
*/

package models

import "time"

/*
	商品服务涉及表格为：VideoType，Video，Membership
	秒杀系统涉及的表格为：FlashSaleProduct,FlashSaleEvent
	为了简化设计，这里将商品服务，秒杀服务集成到一个服务中
*/

// 视频类型
type VideoType struct {
	ID          uint      `json:"-" gorm:"primaryKey"`
	Name        string    `json:"-" gorm:"column:name;type:varchar(100);not null"`
	Description string    `json:"-" gorm:"column:description;type:text"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

// 视频详情
type Video struct {
	ID          uint      `json:"-" gorm:"primaryKey"`
	Title       string    `json:"-" gorm:"column:title;type:varchar(255);not null"`
	TypeID      uint      `json:"-" gorm:"not null"`
	VideoType   VideoType `json:"-" gorm:"foreignKey:TypeID"`
	CoverURL    string    `json:"-" gorm:"column:coverurl;type:varchar(255);not null"`
	PlayURL     string    `json:"-" gorm:"column:playurl;type:varchar(255);not null"`
	Description string    `json:"-" gorm:"column:description;type:text"`
	AccessLevel string    `json:"-" gorm:"column:accessLevel;type:varchar(50);not null"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

// 会员表
type Membership struct {
	ID          uint    `json:"-"  gorm:"primaryKey"`
	Name        string  `json:"-"  gorm:"column:name;type:varchar(100);not null"`
	Price       float64 `json:"-"  gorm:"column:price;not null"`
	Duration    int     `json:"-"  gorm:"column:duration;not null"`
	Description string  `json:"-"  gorm:"column:description;type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (m Membership) TableName() string {
	return "membership"
}

// 秒杀商品表
type FlashSaleEventProduct struct {
	ID                uint      `json:"-"  gorm:"primaryKey"`
	EventID           uint      `json:"-"  gorm:"column:eventid;not null"`           // 加索引
	ProductID         uint      `json:"-"  gorm:"column:productid;not null"`         // 外键，关联到Product表
	OriginalPrice     float64   `json:"-"  gorm:"column:originalprice;not null"`     // 原价
	FlashSalePrice    float64   `json:"-"  gorm:"column:flashsaleprice;not null"`    // 秒杀价
	Quantity          int       `json:"-"  gorm:"column:quantity;not null"`          // 秒杀商品总量
	RemainingQuantity int       `json:"-"  gorm:"column:remainingquantity;not null"` // 剩余秒杀商品数量
	LimitPerUser      int       `json:"-"  gorm:"column:limitperuser;not null"`      // 每用户限购数量
	CreatedAt         time.Time // 创建时间
	UpdatedAt         time.Time // 更新时间
}

func (f FlashSaleEventProduct) TableName() string {
	return "flasheventproduct"
}

// 秒杀活动表
type FlashSaleEvent struct {
	ID        uint      `json:"-"  gorm:"primaryKey"`
	Name      string    `json:"-"  gorm:"column:name;not null"`      // 活动名称
	Condition string    `json:"-"  gorm:"column:condition;not null"` // 活动门槛
	StartTime time.Time `json:"-"  gorm:"column:starttime;not null"` // 活动开始时间
	EndTime   time.Time `json:"-"  gorm:"column:endtime;not null"`   // 活动结束时间
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间
}

func (f FlashSaleEvent) TableName() string {
	return "flashsaleevent"
}
