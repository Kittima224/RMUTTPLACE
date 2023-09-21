package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	StoreId int
	// Images     []PhotoProduct `gorm:"foreignKey:ProductID"`
	Image      string
	Name       string
	Desc       string
	CategoryID int
	Available  int
	Price      int
	Weight     int
	Category   Category
}

// type PhotoProduct struct {
// 	gorm.Model
// 	StoreID   uint
// 	ProductID uint
// 	Image     string
// }
