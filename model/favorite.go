package model

import "gorm.io/gorm"

type Favorite struct {
	gorm.Model
	UserID    uint
	ProductID uint
	Product   Product
}
