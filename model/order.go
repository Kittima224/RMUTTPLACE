package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID   uint
	Products []OrderItem `gorm:"foreignKey:OrderID"`
	Tracking string
	StoreID  uint
	// Price    int
}

type OrderItem struct {
	gorm.Model
	OrderID   uint
	ProductID uint
	// StoreID   uint
	Quantity int
	Product  Product `gorm:"foreignKey:ProductID"`
}
