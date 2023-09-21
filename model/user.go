package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email       string
	UserName    string
	Password    string
	Address     string
	District    string
	SubDistrict string
	Province    string
	Zipcode     string
	Tel         string
	Image       string
}

type Cart struct {
	gorm.Model
	UserID    uint
	ProductID uint
	Quantity  int
	Product   Product `gorm:"foreignKey:ProductID"`
	//Photos    PhotoProduct `gorm:"-"`
}

type Favorite struct {
	gorm.Model
	UserID    uint    `gorm:"foreignKey:UserID"`
	ProductID uint    `gorm:"foreignKey:ProductID"`
	Product   Product `gorm:"foreignKey:ProductID"`
}
