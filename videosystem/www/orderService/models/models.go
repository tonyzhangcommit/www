/*
这里是数据库表设计
*/

package models

import "time"

type Order struct {
	ID            uint    `gorm:"primaryKey"`
	UserID        uint    `gorm:"column:userid;not null"`
	OrderType     string  `gorm:"column:ordertype;type:varchar(100);not null"` // 订单类型
	Status        string  `gorm:"column:status;type:varchar(100);not null"`
	PaymentStatus string  `gorm:"column:paymentstatus;type:varchar(100);not null"`
	TotalAmount   float64 `gorm:"column:totalamount;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type OrderItem struct {
	ID          uint    `gorm:"primaryKey"`
	OrderID     uint    `gorm:"column:orderid;not null"`
	ProductID   uint    `gorm:"column:productid;not null"`
	ProductName string  `gorm:"column:productname;type:varchar(255);not null"` // 商品名称
	Quantity    int     `gorm:"column:quantity;not null"`
	UnitPrice   float64 `gorm:"column:unitprice;not null"`
	Subtotal    float64 `gorm:"column:subtotal;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
