package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID     uint
	Products   []OrderItem `gorm:"foreignKey:OrderID"`
	Tracking   string
	StoreID    uint
	Store      Store
	ShipmentID uint
	Shipment   Shipment
	User       User
}

type OrderItem struct {
	gorm.Model
	OrderID   uint
	ProductID uint
	Quantity  int
	Product   Product `gorm:"foreignKey:ProductID"`
}
