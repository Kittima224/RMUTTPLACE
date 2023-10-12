package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	StoreID     int
	Store       Store
	Image       string
	Name        string
	Description string
	CategoryID  int
	Available   int
	Price       int
	Weight      int
	Category    Category
	Reviews     []Review `gorm:"foreignKey:ProductID"`
	Rating      float32
}

type Review struct {
	gorm.Model
	ProductID int
	UserID    int
	User      User `gorm:"foreignKey:UserID"`
	Comment   string
	Rating    float32
}
