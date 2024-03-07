/*
这里是数据库表设计
*/

package models

import "time"

type Order struct {
	ID            uint    `gorm:"primaryKey"`
	UserID        uint    `gorm:"not null"`
	OrderType     string  `gorm:"type:varchar(100);not null"` // 订单类型
	Status        string  `gorm:"type:varchar(100);not null"`
	PaymentStatus string  `gorm:"type:varchar(100);not null"`
	TotalAmount   float64 `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type OrderItem struct {
	ID          uint    `gorm:"primaryKey"`
	OrderID     uint    `gorm:"not null"`
	ProductID   uint    `gorm:"not null"`
	ProductName string  `gorm:"type:varchar(255);not null"` // 商品名称
	Quantity    int     `gorm:"not null"`
	UnitPrice   float64 `gorm:"not null"`
	Subtotal    float64 `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
