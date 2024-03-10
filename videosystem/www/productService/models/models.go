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
	Description string    `json:"-" gorm:"type:text"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

// 视频详情
type Video struct {
	ID          uint      `json:"-" gorm:"primaryKey"`
	Title       string    `json:"-" gorm:"type:varchar(255);not null"`
	TypeID      uint      `json:"-" gorm:"not null"`
	VideoType   VideoType `json:"-" gorm:"foreignKey:TypeID"`
	CoverURL    string    `json:"-" gorm:"type:varchar(255);not null"`
	PlayURL     string    `json:"-" gorm:"type:varchar(255);not null"`
	Description string    `json:"-" gorm:"type:text"`
	AccessLevel string    `json:"-" gorm:"type:varchar(50);not null"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

// 会员表
type Membership struct {
	ID          uint    `json:"-"  gorm:"primaryKey"`
	Name        string  `json:"-"  gorm:"type:varchar(100);not null"`
	Price       float64 `json:"-"  gorm:"not null"`
	Duration    int     `json:"-"  gorm:"not null"`
	Description string  `json:"-"  gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (m Membership) TableName() string {
	return "membership"
}

// 秒杀商品表
type FlashSaleEventProduct struct {
	ID                uint      `json:"-"  gorm:"primaryKey"`
	EventID           uint      `json:"-"  gorm:"not null"`
	ProductID         uint      `json:"-"  gorm:"not null"` // 外键，关联到Product表
	OriginalPrice     float64   `json:"-"  gorm:"not null"` // 原价
	FlashSalePrice    float64   `json:"-"  gorm:"not null"` // 秒杀价
	Quantity          int       `json:"-"  gorm:"not null"` // 秒杀商品总量
	RemainingQuantity int       `json:"-"  gorm:"not null"` // 剩余秒杀商品数量
	LimitPerUser      int       `json:"-"  gorm:"not null"` // 每用户限购数量
	CreatedAt         time.Time // 创建时间
	UpdatedAt         time.Time // 更新时间
}

// 秒杀活动表
type FlashSaleEvent struct {
	ID        uint      `json:"-"  gorm:"primaryKey"`
	Name      string    // 活动名称
	Condition string    // 活动门槛
	StartTime time.Time `json:"-"  gorm:"not null"` // 活动开始时间
	EndTime   time.Time `json:"-"  gorm:"not null"` // 活动结束时间
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间
}
