package model

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID    uint
	ProductID uint
	Quantity  int
	StoreID   uint
	Product   Product
	Store     Store
}
